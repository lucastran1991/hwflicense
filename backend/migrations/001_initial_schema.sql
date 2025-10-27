-- Customer Master License
CREATE TABLE IF NOT EXISTS cml (
    id TEXT PRIMARY KEY,
    org_id TEXT UNIQUE NOT NULL,
    max_sites INTEGER NOT NULL,
    validity TEXT NOT NULL,
    feature_packs TEXT,
    dev_key_public TEXT NOT NULL,
    prod_key_public TEXT NOT NULL,
    cml_data TEXT NOT NULL,
    signature TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Site Licenses
CREATE TABLE IF NOT EXISTS site_licenses (
    id TEXT PRIMARY KEY,
    site_id TEXT UNIQUE NOT NULL,
    org_id TEXT NOT NULL,
    fingerprint TEXT,
    license_data TEXT NOT NULL,
    signature TEXT NOT NULL,
    issued_at TEXT DEFAULT (datetime('now')),
    last_seen TEXT,
    status TEXT DEFAULT 'active',
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (org_id) REFERENCES cml(org_id)
);

-- Usage Ledger
CREATE TABLE IF NOT EXISTS usage_ledger (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    entry_type TEXT,
    site_id TEXT,
    data TEXT,
    signature TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Usage Statistics
CREATE TABLE IF NOT EXISTS usage_stats (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    period TEXT NOT NULL,
    user_stats TEXT,
    site_stats TEXT,
    total_active_sites INTEGER,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, period)
);

-- Usage Manifests
CREATE TABLE IF NOT EXISTS usage_manifests (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    period TEXT NOT NULL,
    manifest_data TEXT NOT NULL,
    signature TEXT NOT NULL,
    sent_to_astack INTEGER DEFAULT 0,
    sent_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Key Storage (for Hub org keys)
CREATE TABLE IF NOT EXISTS org_keys (
    id TEXT PRIMARY KEY,
    org_id TEXT UNIQUE NOT NULL,
    key_type TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    public_key TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_site_licenses_org_id ON site_licenses(org_id);
CREATE INDEX IF NOT EXISTS idx_site_licenses_status ON site_licenses(status);
CREATE INDEX IF NOT EXISTS idx_usage_ledger_org_id ON usage_ledger(org_id);
CREATE INDEX IF NOT EXISTS idx_usage_stats_org_id ON usage_stats(org_id, period);
CREATE INDEX IF NOT EXISTS idx_usage_manifests_org_id ON usage_manifests(org_id, period);
