# License Metadata Examples

This directory contains example JSON files showing how to structure metadata for different license types when generating licenses via the KMS API.

## Files

### License Type Examples

- **`cml-metadata-example.json`** - Customer Master License (CML) metadata for Hub/Organization level
- **`enterprise-metadata-example.json`** - Enterprise license metadata
- **`site-metadata-prod-example.json`** - Production site license (HWF) metadata
- **`site-metadata-dev-example.json`** - Development site license (Boost - Dev mode) metadata
- **`site-metadata-plant-id-example.json`** - Site license with Plant ID (when Site ID unavailable)
- **`site-metadata-complete-example.json`** - Complete site license with all fingerprint fields
- **`generic-license-metadata-example.json`** - Generic/basic license metadata
- **`trial-license-metadata-example.json`** - Trial license metadata
- **`api-request-example.json`** - Complete API request example with curl command

## Usage

### Using the Metadata in API Request

When calling `POST /licenses/generate`, extract the `metadata` field from these example files:

```bash
# Example: Generate a site license
curl -X POST http://localhost:8080/licenses/generate \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "7ff272a4-1427-4d83-8b72-fb9e1852bf08",
    "license_type": "site",
    "metadata": {
      "site_id": "SITE-2024-001",
      "enterprise_id": "ENT-001",
      "site_name": "Main Manufacturing Plant",
      "mode": "prod",
      "site_type": "hwf",
      "address": "456 Factory Road, Oakland, CA 94601",
      "dns_suffix": "plant1.acme-manufacturing.com",
      "deployment_tag": "hwf-prod-001",
      "max_users": "100",
      "status": "active"
    }
  }'
```

### Extracting Metadata from Example Files

```bash
# Extract metadata from example file and use in request
METADATA=$(cat example/site-metadata-prod-example.json | jq -c '.metadata')

curl -X POST http://localhost:8080/licenses/generate \
  -H "Content-Type: application/json" \
  -d "{
    \"key_id\": \"your-key-id\",
    \"license_type\": \"site\",
    \"metadata\": $METADATA
  }"
```

## Important Notes

1. **All metadata values must be strings** - The current implementation uses `map[string]string`, so:
   - Numbers must be quoted: `"max_users": "100"` (not `100`)
   - Dates must be ISO 8601 strings: `"validity": "2026-12-31T23:59:59Z"`
   - Lists use comma-separated strings: `"features": "feature1,feature2,feature3"`

2. **Fingerprint Fields** - Optional but recommended for site licenses:
   - `address` - Physical address
   - `dns_suffix` - DNS domain suffix
   - `deployment_tag` - Deployment identifier

3. **Key Modes** - For site licenses:
   - `mode: "prod"` - Production mode (required for HWF sites)
   - `mode: "dev"` - Development mode (for Boost sites)

4. **Site Types**:
   - `site_type: "hwf"` - HWF sites (always Prod mode)
   - `site_type: "boost"` - Boost sites (can be Prod or Dev mode)

## License Types

- **Customer Master License (CML)**: Root-level license for Hub/Organization
- **Enterprise License**: Second-level license signed by Hub
- **Site License**: Third-level license signed by Enterprise
- **Generic License**: Basic license for general use
- **Trial License**: Time-limited trial license

## Related Documentation

- See `../requirement.md` for detailed requirements
- See `../project.md` for project structure and API documentation
- See `../kms/README.md` for backend API documentation

