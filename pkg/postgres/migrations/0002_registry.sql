-- haistack-postgres registry catalog schema
--
-- definition_resource and definition_target are global (shared base catalog).
-- registry_install is tenant-scoped enablement/install overlay.

CREATE TABLE IF NOT EXISTS definition_resource (
    canonical_url      TEXT NOT NULL,
    version            TEXT NOT NULL,
    fhir_version       TEXT NOT NULL,
    fhir_resource_type TEXT NOT NULL,
    definition_kind    TEXT NOT NULL,
    name               TEXT NOT NULL,
    status             TEXT NOT NULL,
    package_name       TEXT NOT NULL DEFAULT '',
    package_version    TEXT NOT NULL DEFAULT '',
    module_name        TEXT NOT NULL DEFAULT '',
    json_data          JSONB NOT NULL,
    installed_at       TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (canonical_url, version)
);

CREATE INDEX IF NOT EXISTS idx_definition_resource_kind
    ON definition_resource (definition_kind, fhir_version);

CREATE TABLE IF NOT EXISTS definition_target (
    canonical_url        TEXT NOT NULL,
    version              TEXT NOT NULL,
    target_resource_type TEXT NOT NULL,
    target_role          TEXT NOT NULL,
    PRIMARY KEY (canonical_url, version, target_resource_type, target_role),
    FOREIGN KEY (canonical_url, version)
        REFERENCES definition_resource (canonical_url, version) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_definition_target_lookup
    ON definition_target (target_resource_type, target_role);

CREATE TABLE IF NOT EXISTS registry_install (
    tenant_id            TEXT NOT NULL REFERENCES tenant (id),
    definition_kind      TEXT NOT NULL,
    canonical_url        TEXT NOT NULL,
    version              TEXT NOT NULL,
    target_resource_type TEXT NOT NULL,
    enabled              BOOLEAN NOT NULL DEFAULT true,
    source_module        TEXT NOT NULL DEFAULT '',
    installed_at         TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, definition_kind, canonical_url, version, target_resource_type)
);

CREATE INDEX IF NOT EXISTS idx_registry_install_target
    ON registry_install (tenant_id, target_resource_type, definition_kind, enabled);
