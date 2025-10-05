# Plan — Scribe (Thoth Client)

Owner: Core Services
Language/Runtime: Go 1.23+ (Linux, macOS, Windows), daemon/service with optional tray

This plan is derived from scribe/plan.yaml and expands milestones into actionable, nested tasks with checkboxes. Subtasks are detailed enough to start implementation.

## M0 — Scaffold & Core Loop

- [ ] Cross-platform builds
  - [ ] Configure builds for darwin/arm64, darwin/amd64, linux/amd64, windows/amd64
  - [ ] CI workflow for lint, test, and matrix builds
- [ ] Config + local state
  - [ ] Config loader `~/.thoth/scribe.yaml` with ENV overrides (SCRIBE_)
  - [ ] Initialize SQLite at `~/.thoth/catalog.db`
  - [ ] Define schema tables: files, envelopes, uploads, events, devices
  - [ ] Migrations and vacuum schedule
- [ ] Device enumeration and watcher
  - [ ] Cross-platform volume detection (macOS, Linux, Windows)
  - [ ] Include/exclude patterns (DCIM, Movies, Pictures; exclude tmp/cache)
  - [ ] Periodic rescan (5m) and debounce logic
- [ ] Daemon lifecycle
  - [ ] Graceful shutdown; signal handling
  - [ ] Tray/CLI stub hooks (to be filled in M5)

## M1 — Ingestion & Spooling

- [ ] Media discovery
  - [ ] Detect DCIM/Photos/Movies folders on mounted devices
  - [ ] File type filters and MIME detection
- [ ] Spooler
  - [ ] Fast copy to local spool with fsync
  - [ ] Compute BLAKE3 during copy; record size and hash
  - [ ] Space management: percent-free threshold + eviction policy
- [ ] Heuristics
  - [ ] JPEG-before-RAW ordering
  - [ ] Smallest-file-first when space constrained

## M2 — Upload Engine

- [ ] Protocol client
  - [ ] Reserve: POST `/v1/objects/reserve` (tenant_id, file_name, size, hash)
  - [ ] Chunks: PUT `/v1/objects/{upload_id}/chunks/{n}` with Content-Range
  - [ ] Commit: POST `/v1/objects/{upload_id}/commit`
- [ ] Reliability
  - [ ] Persistent job queue; survive restarts
  - [ ] Retry/backoff `1s..5m` with jitter; max_retries=10
  - [ ] Pause/resume and bandwidth throttle controls
- [ ] Auth
  - [ ] mTLS client certs or token auth based on config

## M3 — Metadata & Catalog

- [ ] Extraction pipeline
  - [ ] EXIF via exiftool/libexif bindings
  - [ ] Audio/video tags via ffprobe bindings
  - [ ] Normalize to envelope JSON; hash envelope and store
- [ ] Local catalog
  - [ ] Index and search by time/device/tags
  - [ ] Sync with server post-commit for verified uploads

## M4 — Replica Confirmation & Cleanup

- [ ] Poll status
  - [ ] GET `/v1/objects/{content_hash}/status` periodically until `r_met`
  - [ ] Respect `r_policy.min_replicas` (>=2)
- [ ] Safe deletion
  - [ ] Write metadata marker to device (tombstone file)
  - [ ] Delete original only after `r_met=true` + ack state
  - [ ] Audit log event and update local catalog state

## M5 — Observability & Packaging

- [ ] Logs/metrics/traces
  - [ ] Structured JSON logs; rotate daily
  - [ ] Prometheus exporter on port 9191
  - [ ] OTLP trace propagation and sampling
- [ ] UI surfaces
  - [ ] System tray (macOS/Windows) with status; CLI status for Linux
  - [ ] Notifications on completion/error
- [ ] Installers/manifests
  - [ ] Linux: deb/rpm + systemd unit
  - [ ] macOS: signed .app bundle + LaunchAgent plist + tray helper
  - [ ] Windows: MSI + service registration + tray
  - [ ] Self-update channel with signature verification

---

## Subsystems

### Watcher
- [ ] OS integration to detect device mount/unmount
- [ ] Scan walker with include/exclude rules
- [ ] Debounce and coalescing of events

### Spool
- [ ] Path resolution `~/.thoth/spool`
- [ ] Disk space monitors; GC policy
- [ ] Integrity: fsync after copy; verify hash

### Uploader
- [ ] Concurrency controls (`concurrent_files`)
- [ ] Chunk size config (`chunk_size`)
- [ ] Resume at chunk N using server state

### Catalog (SQLite)
- [ ] Schema migrations and indices
- [ ] APIs to insert/find files, envelopes, uploads, devices
- [ ] Vacuum interval (24h) task

### R-policy
- [ ] Poll cadence and backoff
- [ ] Device cleanup rules and tombstone markers
- [ ] User-configurable per-device policies

### Security
- [ ] TLS trust store and client cert loading
- [ ] Token auth flows with expiry/rotation

### Observability
- [ ] Counters: files_discovered, uploads_started/failed, files_deleted
- [ ] Gauges: spool_bytes_used, active_uploads
- [ ] Histograms: upload_duration_seconds
- [ ] Notifications: OS/tray and CLI summaries

## Testing
- [ ] Unit: hasher, EXIF parser, retry logic
- [ ] Integration: device mount simulation; end-to-end with local server
- [ ] Performance: 10k-file SD ingest; network drop/resume
- [ ] Security: mTLS handshake; token expiry/rotation

## Packaging & Deployment
- [ ] Linux deb/rpm; systemd unit + config
- [ ] macOS app and LaunchAgent; hardened runtime and signing
- [ ] Windows MSI and service; code signing
- [ ] Auto-updater with signed channel

## Risks & Mitigations
- [ ] Spool exhaustion under burst; enforce low-water marks and eviction
- [ ] Intermittent connectivity; exponential backoff and resume
- [ ] Metadata extraction failures; quarantine and retry policies

## Definition of Done (v1)
- [ ] Cross-platform daemon runs reliably with configurable ingestion
- [ ] Verified uploads with server; R-policy cleanup works safely
- [ ] Usable logs/metrics/traces and minimal UI/tray/CLI surfaces
- [ ] Installers/manifests verified on all platforms

## Backlog Next (post-v1)
- [ ] Delta-sync for renames/moves
- [ ] Compression before upload for large videos
- [ ] LAN peer-to-peer sharing
- [ ] Selective album import (Photos)
- [ ] Scheduled backup windows
- [ ] Energy-aware operation
- [ ] Offline queue for laptops
- [ ] End-to-end encryption with local key escrow
- [ ] Per-device policies (R, compression)
- [ ] Smart duplicate detection across devices
- [ ] iOS shortcut integration
- [ ] CLI diagnostics and repair
- [ ] Multi-server round-robin / geo-failover
- [ ] User analytics dashboard
- [ ] Local HTTP UI for headless
- [ ] Adaptive concurrency by bandwidth
- [ ] Server discovery via mDNS/zeroconf
- [ ] Encrypted external drive spool option
- [ ] Agentless mode (API-triggered)
- [ ] Preview generation in tray UI
