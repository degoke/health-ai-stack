-- haistack-sqlite initial schema
--
-- Future: optional resource_proto_blob column on resource for secondary proto storage.
-- Canonical FHIR JSON in resource.json remains the source of truth.

CREATE TABLE IF NOT EXISTS resource (
    resource_type TEXT NOT NULL,
    id            TEXT NOT NULL,
    version_id    TEXT NOT NULL,
    last_updated  TEXT NOT NULL,
    json          BLOB NOT NULL,
    hash          TEXT,
    created_at    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    PRIMARY KEY (resource_type, id)
);

CREATE TABLE IF NOT EXISTS resource_history (
    rowid         INTEGER PRIMARY KEY AUTOINCREMENT,
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    version_id    TEXT NOT NULL,
    action        TEXT NOT NULL,
    timestamp     TEXT NOT NULL,
    hash          TEXT,
    deleted       INTEGER NOT NULL DEFAULT 0,
    json          BLOB,
    UNIQUE (resource_type, resource_id, version_id)
);

CREATE INDEX IF NOT EXISTS idx_resource_history_resource
    ON resource_history (resource_type, resource_id, rowid);

CREATE TABLE IF NOT EXISTS search_token (
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    field_key     TEXT NOT NULL,
    value         TEXT NOT NULL,
    PRIMARY KEY (resource_type, resource_id, field_key, value)
);

CREATE INDEX IF NOT EXISTS idx_search_token_lookup
    ON search_token (field_key, value);

CREATE TABLE IF NOT EXISTS search_string (
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    field_key     TEXT NOT NULL,
    value         TEXT NOT NULL,
    PRIMARY KEY (resource_type, resource_id, field_key, value)
);

CREATE INDEX IF NOT EXISTS idx_search_string_lookup
    ON search_string (field_key, value);

CREATE TABLE IF NOT EXISTS search_date (
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    field_key     TEXT NOT NULL,
    value         TEXT NOT NULL,
    PRIMARY KEY (resource_type, resource_id, field_key, value)
);

CREATE INDEX IF NOT EXISTS idx_search_date_lookup
    ON search_date (field_key, value);

CREATE TABLE IF NOT EXISTS search_number (
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    field_key     TEXT NOT NULL,
    value         TEXT NOT NULL,
    PRIMARY KEY (resource_type, resource_id, field_key, value)
);

CREATE INDEX IF NOT EXISTS idx_search_number_lookup
    ON search_number (field_key, value);

CREATE TABLE IF NOT EXISTS search_reference (
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    field_key     TEXT NOT NULL,
    value         TEXT NOT NULL,
    PRIMARY KEY (resource_type, resource_id, field_key, value)
);

CREATE INDEX IF NOT EXISTS idx_search_reference_lookup
    ON search_reference (field_key, value);

CREATE TABLE IF NOT EXISTS sync_outbox (
    sequence      INTEGER PRIMARY KEY AUTOINCREMENT,
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    version_id    TEXT NOT NULL,
    action        TEXT NOT NULL,
    timestamp     TEXT NOT NULL,
    hash          TEXT
);

CREATE TABLE IF NOT EXISTS sync_inbox_applied (
    id         TEXT PRIMARY KEY,
    applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_cursor (
    name       TEXT PRIMARY KEY,
    position   TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_conflict (
    id                TEXT PRIMARY KEY,
    resource_type     TEXT NOT NULL,
    resource_id       TEXT NOT NULL,
    local_version_id  TEXT NOT NULL,
    remote_version_id TEXT NOT NULL,
    reason            TEXT NOT NULL,
    created_at        TEXT NOT NULL,
    resolved_at       TEXT
);

CREATE INDEX IF NOT EXISTS idx_sync_conflict_resource
    ON sync_conflict (resource_type, resource_id);

CREATE TABLE IF NOT EXISTS binary_object (
    key          TEXT PRIMARY KEY,
    content_type TEXT,
    size         INTEGER NOT NULL,
    hash         TEXT,
    data         BLOB,
    created_at   TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS module_registry (
    name          TEXT PRIMARY KEY,
    version       TEXT NOT NULL,
    metadata      TEXT,
    registered_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);
