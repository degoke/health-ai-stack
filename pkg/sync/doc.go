// Package sync defines outbox and sync pipeline contracts for haistack-core.
//
// haistack-sync is the event-emission boundary between local resource writes and
// downstream replication, projection, and worker pipelines. It specifies how change
// events are appended without owning cursor management, inbox idempotency, conflict
// reconciliation, or transport protocols.
//
// Callers may alias this package (for example hasync) when they also import the
// standard library sync package.
//
// # Design principles
//
// Append-only change notification:
//
//   - Outbox.Append stores a store.ResourceEvent describing one resource change.
//   - Events mirror history at the notification layer: ResourceType, ID, VersionID,
//     Action (create/update/delete), Timestamp, and content Hash.
//   - store.EventStore implementations (sqlite OutboxStore, postgres EventStore) assign
//     monotonic Sequence values when the caller leaves Sequence unset.
//
// Transactional emission via write sessions:
//
//   - pkg/core appends events inside store.WriteSession so resource, history, index,
//     and event rows commit or roll back together.
//   - WithWriteSession binds a service-level Outbox to a session-scoped EventStore for
//     the duration of one Append call.
//   - EventStoreOutbox adapts a bare store.EventStore to the Outbox interface for
//     non-session wiring and tests.
//
// Optional in core:
//
//   - When pkg/core.ResourceServiceConfig.Outbox is nil, no events are emitted.
//   - When Outbox is non-nil, core emits one event per successful create, update, or
//     delete through WithWriteSession.
//
// Policy stays downstream:
//
//   - Cursor checkpoints use store.CursorStore (not defined here).
//   - Remote apply idempotency (sqlite InboxStore) lives in pkg/sqlite.
//   - Conflict records use store.ConflictStore; reconciliation lives in pkg/core or
//     future sync workers.
//   - Transport (HTTP, MQTT, file drop) is outside this package.
//
// # Outbox interface
//
//	Append(ctx context.Context, event store.ResourceEvent) (store.ResourceEvent, error)
//
// Returns the persisted event, including Sequence when assigned by the backend.
//
// # Helpers
//
// EventStoreOutbox — adapts store.EventStore to Outbox:
//
//	outbox := &sync.EventStoreOutbox{Events: db.OutboxStore()}
//
// SessionOutbox — appends through a WriteSession's EventStore:
//
//	sessionOutbox := &sync.SessionOutbox{Session: session}
//
// WithWriteSession — routes Append transactionally for core writes:
//
//   - When the configured Outbox is *EventStoreOutbox, Append uses the session's
//     EventStore (ignoring the non-transactional Events field on the configured value).
//   - For other Outbox implementations, Append delegates directly to the configured
//     Outbox (for custom hooks, filtering, or test doubles).
//
// Typical core wiring:
//
//	svc, _ := core.NewResourceService(core.ResourceServiceConfig{
//	    …
//	    Outbox: &sync.EventStoreOutbox{Events: tdb.EventStore()},
//	})
//	// core internally calls sync.WithWriteSession(outbox, session).Append(...)
//
// # ResourceEvent fields
//
// Defined in pkg/store/events.go:
//
//   - Sequence     — monotonic backend-assigned id (sqlite AUTOINCREMENT, postgres BIGSERIAL).
//   - ResourceType — FHIR resource type.
//   - ID           — resource id.
//   - VersionID    — version assigned by pkg/core (including tombstone version on delete).
//   - Action       — EventActionCreate, EventActionUpdate, or EventActionDelete.
//   - Timestamp    — write timestamp from core.
//   - Hash         — SHA-256 of canonical JSON at write time (content hash on delete).
//
// # Integration with other packages
//
// pkg/core    — optional Outbox collaborator; emits events on successful writes.
// pkg/store   — ResourceEvent, EventStore, WriteSession, CursorStore contracts.
// pkg/sqlite  — OutboxStore implements EventStore (sync_outbox table).
// pkg/postgres — EventStore implements EventStore (event_log table, global sequence).
//
// # Typical flows
//
// Enable outbox in core (local sqlite):
//
//	svc, _ := core.NewResourceService(core.ResourceServiceConfig{
//	    Resources: db.ResourceStore(),
//	    History:   db.HistoryStore(),
//	    Sessions:  db,
//	    Outbox:    &sync.EventStoreOutbox{},
//	})
//
// Replay events from a cursor (worker layer, not this package):
//
//	events, err := db.OutboxStore().ReadSince(ctx, lastSequence, 100)
//	… process …
//	cursorStore.UpsertCursor(ctx, store.Cursor{Name: "replicator", Position: seq})
//
// Custom outbox hook (must not break transaction unless intentional):
//
//	type auditOutbox struct{ inner sync.Outbox }
//	func (a auditOutbox) Append(ctx context.Context, e store.ResourceEvent) (store.ResourceEvent, error) {
//	    // side effect
//	    return sync.WithWriteSession(a.inner, sessionFromContext(ctx)).Append(ctx, e)
//	}
//
// # File layout
//
//   - doc.go            — package documentation (this file)
//   - outbox.go         — Outbox interface, EventStoreOutbox
//   - session_outbox.go — SessionOutbox, WithWriteSession
//
// # Testing
//
// outbox_test.go verifies WithWriteSession routes EventStoreOutbox to session EventStore
// and delegates custom Outbox implementations directly.
//
// pkg/core tests cover outbox failure rollback and delete tombstone versionId on events.
//
// # Out of scope (MVP)
//
// Inbox idempotency interface (sqlite InboxStore exists without pkg/store contract),
// cursor management, conflict reconciliation workers, transport adapters, batch replay,
// deduplication policy, and cross-tenant fan-out. This package defines event append
// contracts and session binding helpers only.
package sync
