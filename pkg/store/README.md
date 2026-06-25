# haistack-store (`pkg/store`)

Storage contract layer for the health-ai-stack monorepo.

## What it does

Think of **`pkg/types`** as “how we represent a FHIR resource in memory.”  
**`pkg/store`** is “how we *save*, *find*, and *track changes* to those resources — without caring whether the database is SQLite, Postgres, or something else.”

This package **does not talk to a database directly**. It defines **interfaces** (contracts) such as:

- save, read, update, and delete a resource
- keep version history
- index data for search
- record change events for sync
- store files, audit logs, background jobs, and more

Actual database code lives in separate packages (`pkg/sqlite`, `pkg/postgres`). Application code depends on **`store.ResourceStore`**, not on a specific backend. That makes it easy to swap implementations or test with in-memory fakes.

## What it does not do

- Parse FHIR JSON — use `pkg/types`
- Assign version numbers or enforce business rules — use `pkg/core` or your app layer
- Parse FHIR search queries — use `pkg/search`
- Run SQL or manage schemas — use `pkg/sqlite` or `pkg/postgres`

`pkg/store` is the **middle layer**: it defines *what* can be persisted and *how* to call it. Backends decide *where* data actually goes.

## When to use it

- When writing or reading FHIR resources through a stable API
- When you need history, search indexes, or change events alongside current state
- When building sync, audit, or background-job features on top of shared storage contracts
- When you want tests to use fake stores instead of a real database

## Mental model

When a Patient is saved, several things usually happen:

1. **Current state** — store the latest Patient (`ResourceStore`)
2. **History** — remember this was version 3 (`HistoryStore`)
3. **Search index** — index `"family=Doe"` so search works (`SearchStore`)
4. **Change event** — emit “Patient/pat-1 was updated” for sync (`EventStore`)

`pkg/store` defines separate interfaces for each of those jobs, plus helpers for sync cursors, conflicts, files, audit trails, and jobs.

## Usage

**Open a backend and get a store:**

```go
import (
    "github.com/degoke/health-ai-stack/pkg/sqlite"
)

db, err := sqlite.Open(sqlite.Options{Path: "data.db"})
resources := db.ResourceStore() // satisfies store.ResourceStore
```

On a server you might use `pkg/postgres` instead — same interfaces, different backend.

**Save and read a resource:**

```go
// envelope comes from pkg/types
err := resources.Create(ctx, envelope)

patient, err := resources.Read(ctx, "Patient", "pat-1")
```

**Compose a full write (typical pattern):**

Higher-level code (or a `WriteSession`) often does all of this together:

```go
resources.Create(ctx, envelope)

history.AppendVersion(ctx, store.ResourceVersion{
    ResourceType: envelope.ResourceType,
    ID:           envelope.ID,
    VersionID:    envelope.VersionID,
    Action:       store.VersionActionCreate,
    Timestamp:    time.Now(),
    Resource:     envelope,
    Hash:         envelope.Hash,
})

events.Append(ctx, store.ResourceEvent{
    ResourceType: envelope.ResourceType,
    ID:           envelope.ID,
    VersionID:    envelope.VersionID,
    Action:       store.EventActionCreate,
    Timestamp:    time.Now(),
    Hash:         envelope.Hash,
})

search.Index(ctx, store.SearchIndexEntry{
    ResourceType: "Patient",
    ID:           "pat-1",
    Fields:       map[string]string{"family": "Doe"},
})
```

On delete: remove from `ResourceStore`, append a tombstone to history, emit a delete event, and remove from the search index.

**Sync and background work:**

- `EventStore` + `CursorStore` — replay changes and resume where you left off
- `ConflictStore` — record “local vs remote version mismatch”
- `JobStore` — queue work like reindexing
- `AuditStore` — record who did what for compliance

## Interface overview

| Interface | Purpose |
|-----------|---------|
| `ResourceStore` | Latest copy of each resource |
| `HistoryStore` | All past versions |
| `SearchStore` | Fast lookup data (not full search parsing) |
| `EventStore` | Change feed for sync and replay |
| `CursorStore` | “Where did the sync worker stop?” |
| `ConflictStore` | Record sync/write conflicts |
| `BinaryStore` | Small inline binary payloads |
| `BlobStore` | Larger blobs or external storage references |
| `MaterializedViewStore` | Read-optimized projections |
| `AnalyticsStore` | Operational/product analytics events |
| `AuditStore` | Security/compliance audit trail |
| `JobStore` | Background task queue |
| `ModuleStore` | Registered plugins/extensions |
| `IDRegistryStore` | Authoritative ID registration |
| `WriteSessionProvider` | Atomic write across resource, history, search, and events |

## Where it fits

| Layer | Role |
|-------|------|
| **types** | Canonical JSON and `ResourceEnvelope` |
| **store** | Persistence contracts (this package) |
| **sqlite** | Embedded local implementation |
| **postgres** | Tenant-scoped edge/cloud implementation |
| **core** | Business rules and write pipelines |
| **search** | FHIR search parsing (produces index entries for `SearchStore`) |

## MVP limits

- Contracts and record types only — no database adapters in this package
- No FHIR search parsing, version assignment policy, or conflict reconciliation
- No in-memory production adapter yet (`pkg/store/memory` is future work)
- `BinaryStore`, `BlobStore`, and operational stores are defined here; not every backend implements every interface

See [doc.go](./doc.go) for the full API, design principles, typical flows, and file layout.
