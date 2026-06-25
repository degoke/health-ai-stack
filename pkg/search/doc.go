// Package search defines FHIR search indexing contracts for haistack-core.
//
// haistack-search is the indexing boundary between persisted FHIR resources and
// pkg/store search tables. It specifies how resources are translated into
// store.SearchIndexEntry values without implementing SearchStore or parsing raw
// FHIR search query strings at runtime.
//
// # Design principles
//
// Index extraction, not query execution:
//
//   - Indexer.Build produces []store.SearchIndexEntry from a types.ResourceEnvelope.
//   - store.SearchStore (pkg/sqlite, pkg/postgres) persists entries and answers lookups.
//   - FHIR search parameter parsing, chained searches, _include, _revinclude, and
//     _filter evaluation remain outside this package.
//
// Post-persistence indexing in core:
//
//   - pkg/core invokes Indexer after resource and history writes inside a WriteSession.
//   - On create/update: core removes prior index rows, then indexes new entries.
//   - On delete: core removes index rows when Indexer is configured (even though Build
//     is not called for deletes).
//   - Indexer failures abort the WriteSession (no partial commit).
//
// Typed field keys:
//
//   - SearchIndexEntry.Fields maps logical keys to string values.
//   - Backends route keys to typed tables using prefixes (see pkg/sqlite and
//     pkg/postgres documentation):
//     token.*, string.*, date.*, number.*, reference.* / ref.*
//   - Keys without a prefix default to string indexing in sqlite/postgres adapters.
//
// Optional collaborator:
//
//   - pkg/core treats a nil Indexer as a no-op (resources persist without search rows).
//
// # Indexer interface
//
//	Build(ctx context.Context, resource *types.ResourceEnvelope) ([]store.SearchIndexEntry, error)
//
// Implementations should:
//
//   - Return one or more SearchIndexEntry values with Fields populated for lookup.
//   - Set ResourceType and ID on entries when known; core overwrites them with the
//     envelope's canonical type and id before indexing.
//   - Return nil or an empty slice when the resource has no indexable fields.
//   - Return an error to abort the write when extraction fails.
//
// SearchIndexEntry shape (from pkg/store):
//
//   - ResourceType — FHIR resource type (for example "Patient").
//   - ID           — resource id.
//   - Fields       — map of index key → string value (for example "string.family" → "Doe").
//
// # Integration with other packages
//
// pkg/core   — invokes Indexer on successful create/update inside WriteSession.
// pkg/store  — SearchStore persists and queries SearchIndexEntry values.
// pkg/types  — ResourceEnvelope.JSON is the typical extraction source.
// pkg/sqlite — typed search_* tables back SearchStore on embedded nodes.
// pkg/postgres — tenant-scoped search_* tables back SearchStore on server nodes.
//
// # Typical flows
//
// Indexer wired into core:
//
//	svc, _ := core.NewResourceService(core.ResourceServiceConfig{
//	    Resources: db.ResourceStore(),
//	    History:   db.HistoryStore(),
//	    Sessions:  db,
//	    Indexer:   patientIndexer{},
//	})
//
// Lookup after indexing (application or search service layer):
//
//	ids, err := db.SearchStore().Lookup(ctx, "string.family", "Doe")
//
// # File layout
//
//   - doc.go    — package documentation (this file)
//   - search.go — Indexer interface
//
// # Out of scope (MVP)
//
// FHIR search parameter parser, composite/chained token handling, _content full-text,
// _filter, sorting, paging, include/revinclude resolution, OpenSearch adapters, and
// SearchStore implementations. Query planning and HTTP search endpoints belong in future
// packages built on top of store.SearchStore and Indexer output.
package search
