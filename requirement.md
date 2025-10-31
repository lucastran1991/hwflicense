# Key Management and Licensing Process

## 1. Key Generation and Storage

- Key generation occurs on the A-Stack License Server side.
- Keys have expiration dates and must be refreshed via the A-Stack License Server.
- During key refresh, the license server must ensure the prior key is valid.
- There are three key levels: Hub key, Enterprise key, and Site key.
- Site keys are of two kinds: Dev Key and Prod Key.
- All keys should be securely stored using the standard AWS key store for encryption, management, and access control.

## 2. Provisioning

- Issue a “Customer Master License” (CML) to the Hub, including:
  - `org_id`, `max_enterprise`, `max_sites`, `max_users`, `validity`, and optional feature packs.
- The CML is signed by you; the Hub validates it with your public key.

## 3. Enterprise Creation

- For each Enterprise, the operator uses the Hub UI/CLI to enter:
  - Enterprise Name, Enterprise ID, and optional fingerprint fields (address, DNS suffix, deployment tag).
- The Hub mints an Enterprise license (`enterprise.lic`), signed with the Hub’s org key.
- This signature chains back to the Root key, ensuring trust.

## 4. Site Creation

- For each new site, the operator uses Hub UI/CLI to enter:
  - `site_id` and optional fingerprint fields (address, DNS suffix, deployment tag).
- The Enterprise mints a `site.lic` (sub-license), signed with the Enterprise’s key, chaining back to your root.
- If a site is created in HWF but the Site ID is not yet available, the system can use the internal Plant ID to register the site on the license server.

## 5. Runtime Behavior

- Each site node uses `site.lic`. Your app verifies:
  - Chain of trust (Root → Org Hub → Enterprise license → Site license).
  - Constraints (features, expiry, fingerprint).

## 6. Counting and Usage Manifest

- The Hub maintains the ground truth—a signed ledger of all active sites.
  - Each record includes: `site_id`, `issued_at`, and the `last_seen` heartbeat timestamp.
- Monthly, the Hub emits a Usage Manifest (JSON + signature) summarizing:
  - Total active site count and details.
  - User counts:
    - HWF Admins (including list of HWF admin users)
    - Enterprise Admins
    - Enterprise Users
    - Plant Users
    - Demo Users
  - Site summaries by Enterprise:
    - Enterprise Name (Enterprise ID)
      - Total number of sites under that Enterprise
      - Individual counts: Boost sites, HWF sites, Dev mode sites, Prod mode sites, Active, Basic, and Commissioning plants
- Customers can email or upload this file to you monthly (configurable).

## 7. Site License Constraints

- Each site can only have one active key (Prod or Dev mode).
- HWF sites are always in Prod mode.
- Boost sites can operate in Prod or Dev mode.
- When a Boost site becomes an HWF site, its key must be updated to Prod mode.
- Whenever a site is created, updated, or deleted, the license server must be updated.
- If duplicate site IDs exist (e.g., commissioning/active or basic), each is treated as a distinct site and must have its own unique site key.

## 8. A-Stack License Server – Core REST APIs

1. Create - Create site/enterprise
2. Get the site API key
3. Update - Update site/enterprise or API key mode
4. Delete - Delete site/enterprise
5. Refresh the Key - Ensure prior key is valid during refresh
6. Get Aggregate Stats
7. Validate Key (Token Check)
   - Verify if the site API key/token from HWF is authentic, active, and valid.
8. Alert/Warning When Key Is Invalid
   - When a key is invalid or expired, HWF sends a message to the central system.
   - A log (e.g., `license_event_log` or `site_key_status`) tracks invalid/expired key events.
   - HWF stops refreshing data from Insight for sites with invalid keys.
   - HWF displays a message, and the plant-side page for that site is blank:
     > “Key for site X is invalid — data for this site is disabled until renewed.”
