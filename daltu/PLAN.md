# Plan — Daltu (MediaBackup Server)

Owner: Core Services
Language/Runtime: Go 1.23+ (Linux-only target), systemd service

This plan is derived from daltu/plan.yaml and breaks work into actionable tasks with checkboxes. Tasks are grouped by milestone and subsystem. Subtasks are detailed enough to begin implementation.

## M0 — Scaffold & Core Wiring

- [ ] Go module, build, and CI
  - [ ] Ensure `go.mod` and reproducible builds for linux/amd64 and linux/arm64
  - [ ] Add CI workflow: lint (staticcheck/golangci-lint), test, and build artifacts
  - [ ] Add Makefile targets: build, test, lint, package
- [ ] Configuration system
  - [ ] Define config struct covering keys in plan.yaml (listen, auth, storage, replication, acks, database, quotas, housekeeping, observability, security)
  - [ ] Load from file `/etc/mediabackup/server.yaml` with ENV overrides (prefix MBK_)
  - [ ] Add validation (required fields, sane defaults, enums)
  - [ ] Provide `--config` flag and dump-effective-config endpoint for debugging
- [ ] Process management
  - [ ] Add systemd unit file (from plan.yaml) and packaging script to install
  - [ ] Graceful shutdown handling and readiness/health probes
- [ ] Observability bootstrap
  - [ ] Structured JSON logger with request_id/tenant_id fields
  - [ ] Basic Prometheus metrics endpoint `/v1/metrics`
  - [ ] OTLP tracing wiring with sampling config

## M1 — Ingest API (local durability)

- [ ] HTTP/JSON API surface selection and scaffolding
  - [ ] Choose HTTP/JSON for v1; reserve gRPC for future
  - [ ] Versioned routes under `/v1`
  - [ ] Middleware: auth, logging, tracing, panic recovery, rate limiting
- [ ] Reserve upload
  - [ ] Endpoint: POST `/v1/objects/reserve`
  - [ ] Request validation (tenant_id, content_hash, size, idempotency_key)
  - [ ] Issue `upload_id`, record in DB `uploads` table with TTL
  - [ ] If object exists: return already_present=true
- [ ] Chunked/resumable upload
  - [ ] Endpoint: PUT `/v1/objects/{upload_id}/chunks/{n}`
  - [ ] Accept Content-Range; verify per-chunk BLAKE3
  - [ ] Persist to staging: `{staging_root}/{upload_id}/{chunk_idx}` and fsync
  - [ ] Update `uploads.received_bytes`, return `next_chunk`
- [ ] Commit upload
  - [ ] Endpoint: POST `/v1/objects/{upload_id}/commit`
  - [ ] Stitch/verify whole-object BLAKE3 and size
  - [ ] Move to CAS path `{cas_root}/{h0}{h1}/{h2}{h3}/{hash}` and fsync
  - [ ] Create `objects` + `envelopes` rows; set `replicas_current`
  - [ ] Ack policy: reply `ack_policy_state` (local_persisted|r_met)
- [ ] Get object status
  - [ ] Endpoint: GET `/v1/objects/{content_hash}/status` with replica locations
  - [ ] Response: replicas list, replicas_current, r_met

## M2 — Replication Controller

- [ ] Targets abstraction
  - [ ] Define interface for `local|s3|nfs|custom` targets
  - [ ] Implement local filesystem copy with fsync and checksum verify
  - [ ] Implement S3 multipart with SSE-KMS option; verify vs BLAKE3
  - [ ] Implement NFS copy with verify-after-write and warnings
- [ ] Scheduler
  - [ ] Scan for objects with `replicas_current < R`
  - [ ] Respect `max_inflight` and failure-domain diversity
  - [ ] Backoff and circuit breaker per target
- [ ] Verification and state updates
  - [ ] Upon successful copy + verify, add `replicas` row, update counts
  - [ ] Emit Event `R_MET` when target R reached (for polling clients)
- [ ] Background jobs
  - [ ] Repair loop to heal missing/unverified replicas
  - [ ] Periodic audit sample full verify

## M3 — Catalog & Queries

- [ ] SQL schema and migrations
  - [ ] Define migrations for `objects`, `envelopes`, `replicas`, `uploads`, `events`
  - [ ] Add indexes (GIN on `envelopes.raw_json`, ts desc on events)
  - [ ] Provide dev SQLite and prod PostgreSQL support
- [ ] Search API
  - [ ] Endpoint: GET `/v1/search` with query params and pagination
  - [ ] Implement simple filter DSL or param-based filters (time, device, tags)

## M4 — Ops, Security & Observability

- [ ] AuthN/Z
  - [ ] Support `mtls|token` modes from config
  - [ ] Validate OIDC tokens with issuer list; tenant scoping via claims
  - [ ] Enforce per-tenant quotas; 429 with Retry-After
- [ ] TLS/mTLS
  - [ ] Load TLS cert/key and CA bundle; require TLS 1.3 strong ciphers
  - [ ] Optional client cert requirement
- [ ] Metrics & Tracing
  - [ ] Counters, gauges, histograms per plan.yaml metrics list
  - [ ] Trace propagation across ingest → storage → replication
- [ ] Runbooks & backup/restore
  - [ ] Document backup/restore of DB and CAS
  - [ ] Fire-drill validation checklist

---

## Subsystem Task Breakdown

### Content-Addressed Storage (CAS)
- [ ] Define path fanout (2 levels) and implement helper funcs
- [ ] Implement atomic move from staging to CAS with fsync
- [ ] Quarantine on mismatch: `/srv/mediabackup/quarantine/{hash}`

### Database
- [ ] Migration tool integration (goose or migrate)
- [ ] Connection pool tuning; timeouts; context plumbing
- [ ] Views: `object_replica_status`

### Idempotency & Resumability
- [ ] Enforce `Idempotency-Key` + hash dedupe
- [ ] Resume from chunk N; expiration and cleanup of stale reservations

### Multi-tenancy
- [ ] `tenant_id` everywhere; scoped quotas
- [ ] Admin endpoints for quota and tenant status

### Observability
- [ ] Logging redaction rules for sensitive fields
- [ ] Alerts: disk free <15%, backlog threshold, checksum mismatch rate

### Housekeeping Jobs
- [ ] GC_Staging (15m)
- [ ] Verify_Replicas (24h)
- [ ] Heal_Replicas (1h)
- [ ] Reindex_Catalog (6h)
- [ ] Quota_Enforcement (10m)

### Testing
- [ ] Unit: chunk hasher, CAS paths, validators
- [ ] Integration: reserve→chunks→commit→replicate→verify
- [ ] Chaos: target outages, corrupt replica, recovery paths
- [ ] Performance: sustained ingest with N clients
- [ ] Security: mTLS, token validation, authZ scopes

### Deployment & Packaging
- [ ] Static binary build and Debian package
- [ ] systemd unit install script; env files; limits
- [ ] Blue/green or canary rollout scripts

### Risks & Mitigations
- [ ] Low disk watermarks halt commits; envelopes-only mode
- [ ] Backpressure and retries on target outage
- [ ] Quarantine mismatches; manual review tooling

## Definition of Done (v1)
- [ ] Clients ingest at scale; resume on failure; poll R until met
- [ ] Integrity enforced end-to-end; envelopes stored; replicas managed
- [ ] Operators have metrics/logs/traces/runbooks
- [ ] Disaster drill proves restore from replicas and DB snapshots

## Backlog Next (post-v1)
- [ ] Erasure coding (Reed–Solomon)
- [ ] WORM retention and legal hold
- [ ] S3 Glacier/Deep Archive tiering
- [ ] Admin web UI
- [ ] Geo-replication with residency policies
- [ ] Lifecycle policies (TTL/ILM)
- [ ] Hash agility migrations
- [ ] End-to-end encryption with KMS/HSM and client-held keys
- [ ] Pluggable metadata extraction workers
- [ ] Near-duplicate detection (pHash/aHash/DCT)
- [ ] Restore/export tooling
- [ ] DB sharding/replicas
- [ ] Event streaming (Kafka/NATS/Webhooks)
- [ ] Policy engine (OPA/Rego)
- [ ] Cross-target consistency checks and quorum acks
- [ ] Offline ingestion bundles (signed manifests)
- [ ] Bandwidth/QoS scheduler
- [ ] Native tus.io endpoints
- [ ] Periodic scrubbing/bit-rot detection
- [ ] S3-compatible facade over CAS
