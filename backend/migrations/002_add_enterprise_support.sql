-- Add enterprise support and new fields from Q&A meeting

-- Enterprises table - enterprise-level keys
CREATE TABLE IF NOT EXISTS enterprises (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    org_id TEXT NOT NULL,
    enterprise_key TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Add new columns to site_licenses table (only if they don't exist)
-- We'll use a separate check since SQLite doesn't support IF NOT EXISTS for ALTER TABLE
-- This migration attempts to add columns, errors will be ignored if columns exist

-- Add columns (ignore errors if they already exist)
-- In production, we handle this by catching the error in the migration runner
ALTER TABLE site_licenses ADD COLUMN key_type TEXT DEFAULT 'production';
ALTER TABLE site_licenses ADD COLUMN expires_at TEXT;
ALTER TABLE site_licenses ADD COLUMN enterprise_id TEXT;
ALTER TABLE site_licenses ADD COLUMN last_refreshed TEXT;

-- Create indexes (these use IF NOT EXISTS and will be skipped if already exist)
CREATE INDEX IF NOT EXISTS idx_site_licenses_enterprise ON site_licenses(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_site_licenses_key_type ON site_licenses(key_type);
CREATE INDEX IF NOT EXISTS idx_site_licenses_expires ON site_licenses(expires_at);

