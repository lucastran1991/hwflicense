-- License Server Tables for Hub Integration
-- These tables support the 7 license-server APIs merged into hub backend

-- Enterprises table - enterprise-level keys
CREATE TABLE IF NOT EXISTS enterprises (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    org_id TEXT NOT NULL,
    enterprise_key TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Site keys table - site-level keys with type (dev/prod), expiration, status
CREATE TABLE IF NOT EXISTS site_keys (
    id TEXT PRIMARY KEY,
    site_id TEXT UNIQUE NOT NULL,
    enterprise_id TEXT NOT NULL,
    key_type TEXT NOT NULL DEFAULT 'production',
    key_value TEXT NOT NULL,
    issued_at TEXT DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    last_validated TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (enterprise_id) REFERENCES enterprises(id)
);

-- Key refresh log - audit trail for monthly refreshes
CREATE TABLE IF NOT EXISTS key_refresh_log (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL,
    old_key TEXT NOT NULL,
    new_key TEXT NOT NULL,
    refreshed_at TEXT DEFAULT (datetime('now')),
    reason TEXT,
    FOREIGN KEY (site_id) REFERENCES site_keys(site_id)
);

-- Quarterly stats - aggregated stats per quarter
CREATE TABLE IF NOT EXISTS quarterly_stats (
    id TEXT PRIMARY KEY,
    period TEXT NOT NULL,
    production_sites INTEGER DEFAULT 0,
    dev_sites INTEGER DEFAULT 0,
    user_counts TEXT,
    enterprise_breakdown TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(period)
);

-- Validation cache - token cache for validation
CREATE TABLE IF NOT EXISTS validation_cache (
    id TEXT PRIMARY KEY,
    site_id TEXT UNIQUE NOT NULL,
    token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (site_id) REFERENCES site_keys(site_id)
);

-- Alerts - invalid key alerts
CREATE TABLE IF NOT EXISTS alerts (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL,
    alert_type TEXT NOT NULL,
    message TEXT NOT NULL,
    alert_timestamp TEXT NOT NULL,
    sent_to_astack INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (site_id) REFERENCES site_keys(site_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_site_keys_enterprise ON site_keys(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_site_keys_status ON site_keys(status);
CREATE INDEX IF NOT EXISTS idx_site_keys_expires ON site_keys(expires_at);
CREATE INDEX IF NOT EXISTS idx_key_refresh_site ON key_refresh_log(site_id);
CREATE INDEX IF NOT EXISTS idx_validation_cache_site ON validation_cache(site_id);
CREATE INDEX IF NOT EXISTS idx_alerts_site ON alerts(site_id);
CREATE INDEX IF NOT EXISTS idx_quarterly_stats_period ON quarterly_stats(period);

