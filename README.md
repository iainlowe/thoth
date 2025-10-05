# Thoth Suite — README

## Overview

**Thoth** is an end-to-end, privacy-preserving media backup and archival system for individuals and small organizations.  
It consists of three tightly integrated components:

- **Daltu** — the **server**: a fault-tolerant, content-addressed storage and replication service.  
- **Scribe** — the **client daemon**: a cross-platform background agent that detects and uploads media.  
- **Stylus** — the **Apple app**: a native iOS/iPadOS/macOS front-end with Photos integration.

Together they form a complete ecosystem where devices capture, store, and synchronize media safely, verifying every byte from lens to archive.

---

## ✳️ Thoth Philosophy

Thoth is designed around three principles:

1. **Integrity** — every object is verified, checksummed, and traceable.  
2. **Resilience** — replication across independent failure domains before deletion.  
3. **Transparency** — complete audit trails, from capture metadata to server replicas.

It acts as a *scribe and library of truth* for your photos, videos, and creative work—built with verifiability, not trust, at its core.

---

## 🧱 Components

### **Daltu – The Server**
**Language:** Go • **Target:** Linux • **Deployment:** systemd / container

Daltu is a high-integrity, content-addressed storage backend.  
It receives uploads from Scribe or Stylus clients, validates their checksums, and replicates objects to durable targets (local disks, S3, NFS, etc.) until the configured *R-value* (replica count) is met.

**Core Features**
- Chunked and resumable ingest API (`/v1/objects/...`)
- Content-addressed storage (BLAKE3 hash fan-out)
- Configurable replica controller (local + cloud backends)
- Envelope catalog (EXIF/IPTC metadata, camera data, GPS)
- SQL catalog with search & query APIs
- Background verification, healing, and replication repair
- mTLS or OIDC-token authentication
- JSON logs, Prometheus metrics, OTLP tracing

**Durability Model**
- Upload to staging → checksum → fsync → CAS commit  
- Replication controller ensures *R* replicas across domains  
- Clients only delete source media once the server confirms `r_met = true`

---

### **Scribe – The Client**
**Language:** Go • **Platforms:** Linux, macOS, Windows  
**Runtime:** Background daemon / tray service

Scribe is the local agent that **detects devices, copies media, and synchronizes** them with Daltu.  
It operates automatically when you plug in an SD card, camera, or phone.  
When space allows, it copies all media to a local *spool* first for speed, then streams verified chunks to Daltu.

**Core Features**
- Cross-platform filesystem watchers for removable devices
- DCIM auto-detection and “JPEG-before-RAW” heuristic
- Fast local spool (optional smallest-file-first fallback when space-limited)
- Envelope extraction (EXIF, ID3, etc.) via exiftool/ffprobe bindings
- Local SQLite catalog with audit history
- R-policy enforcement — keeps originals until required replicas confirmed
- Configurable bandwidth throttling, retry/backoff, and chunk size
- JSON logs, Prometheus metrics, and local CLI or tray UI

**Deployment**
- Linux: systemd service  
- macOS: launchd agent + tray UI  
- Windows: service + taskbar icon  

Scribe operates autonomously: plug in media, it syncs; unplug, it sleeps.

---

### **Stylus – The App**
**Language:** SwiftUI • **Platforms:** iOS, iPadOS, macOS  
**Integration:** Apple Photos, Core Data, Background URLSession

Stylus provides a native Apple experience for **selective backup and verification** of your Photos library.  
It connects directly to Daltu using the same authenticated protocol as Scribe, preserving on-device privacy controls.

**Core Features**
- Deep Photos integration (albums, Live Photos, RAW, ProRes)
- Background uploads via `BGProcessingTask` and `URLSession`
- Rich envelope metadata (EXIF/IPTC/GPS/video tracks)
- Replica-status indicators (shows when R-value met)
- Per-album privacy and export profiles
- Optional restore/export: search server catalog, download, re-import to Photos
- iCloud-optimized asset handling (auto-fetch on demand)
- Privacy redaction: GPS fuzzing, face stripping, hashed Wi-Fi SSIDs
- Configurable export “profiles” (Private / Standard / Rich Archive)

**Design Philosophy**
Stylus doesn’t just back up your photos—it helps you verify their existence, context, and authenticity, maintaining the provenance chain from capture to archive.

---

## 🔗 Architecture Overview

```
 [Stylus]     [Scribe]              [Daltu]
  iOS/macOS    Daemon (Go)           Server (Go)
     │              │                    │
     │   HTTPS/mTLS │                    │
     ├──────────────▶ Ingest API         │
     │              │                    ▼
     │              │          Content-Addressed Storage
     │              │          + Replica Controller
     │              │                    │
     │              ◀────────── Status / R-checks
     │
     ▼
User notified when safe (R replicas met)
```

---

## 🔒 Data Model & Integrity

- **Object Identity** — BLAKE3 hash of canonical file bytes  
- **Envelope** — JSON metadata with camera, EXIF, GPS, and policy info  
- **Replica Record** — verified copies with timestamp & failure domain  
- **Integrity Path** — local → staging → CAS → replica → R_met  
- **Verification** — every hop rehashes & fsyncs before acknowledgment  

This structure guarantees end-to-end verifiability, even across mixed storage backends.

---

## 🧩 Configurability

| Layer | Key Controls | Examples |
|-------|---------------|-----------|
| **Daltu** | `R`, replication targets, ack policy | e.g. 3 replicas (local + NAS + S3) |
| **Scribe** | spool size, retry policy, file filters | 20 GB spool, 8 MB chunks, JPEG-before-RAW |
| **Stylus** | privacy/export profile, album filters | “Favorites only”, “No GPS”, “Wi-Fi only” |

All components share a common configuration schema with YAML/JSON or GUI equivalents.

---

## 🧠 Data Flow Summary

1. **Detection:** Scribe or Stylus detects new media sources.  
2. **Ingestion:** Metadata envelopes are generated locally.  
3. **Verification:** Hash and size computed before upload.  
4. **Transfer:** Chunked uploads to Daltu with retry & resume.  
5. **Replication:** Daltu replicates to durable targets.  
6. **Acknowledgment:** Clients poll until `R_met`.  
7. **Cleanup:** Originals deleted only after verified redundancy.  
8. **Cataloging:** Metadata searchable via Daltu’s API or Stylus UI.

---

## 🛠 Developer & Operator Notes

- **APIs:** REST/JSON with optional gRPC bindings.  
- **Database:** PostgreSQL or SQLite (dev).  
- **Storage:** Local CAS on ext4/XFS + optional S3/NFS backends.  
- **Observability:** Prometheus metrics, OTLP tracing, structured JSON logs.  
- **Deployment:** Single-binary install; minimal dependencies.  
- **Security:** TLS 1.3, mTLS optional, per-tenant tokens, content hashing.

---

## 📦 Example Directory Layout

```
/srv/daltu/
  ├── cas/
  ├── staging/
  ├── replicas/
  ├── quarantine/
  └── db/

~/Library/Application Support/Scribe/
  ├── spool/
  ├── catalog.db
  └── logs/

~/Library/Containers/com.thoth.Stylus/
  ├── CoreData.sqlite
  ├── config.plist
  └── logs/
```

---

## 🌐 Vision Roadmap

- Erasure coding for efficient multi-replica storage  
- Object authenticity (C2PA / content signatures)  
- Family and organization sharing with policy-based access  
- Intelligent deduplication and delta compression  
- Search by metadata, tags, and visual similarity  
- Offline ingestion bundles (“sneakernet mode”)  
- Mobile-to-mobile relay over local network  

---

## 📖 Summary

> **Thoth** is a holistic, verifiable memory system for the modern era.  
> **Daltu** preserves; **Scribe** transcribes; **Stylus** curates.  
> Together, they ensure that every captured moment—whether from camera, phone, or archive—remains intact, authenticated, and retrievable long after devices and clouds have changed.

---

**License:** Apache 2.0 (provisional)  
**Maintainers:** Core Services / Thoth Team  
**Keywords:** backup, replication, authenticity, photo library, agentic storage, provenance
