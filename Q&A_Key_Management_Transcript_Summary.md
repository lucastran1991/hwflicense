# Q&A Key Management - Transcript Summary

**Date:** October 27, 2025, 3:45 PM  
**Participants:** Alok Batra, Kartik Shah, Nancy Tran, An Nguyen Thanh

---

## 📋 Background & Requirements

### Business Context
- **Problem:** A-Stack's access to Veolia's environment will be removed (within ~1 year) as Veolia takes over management
- **Challenge:** Need to track how many sites are active to determine licensing costs
- **Solution:** Implement a licensing module to track and enforce site usage

### Key Business Requirements

#### 1. **Site-based Licensing**
- Charging based on number of sites for HWF and Boost
- Quarterly licensing calculation (not continuous)
- Track active sites without direct environment access

#### 2. **Configuration vs Production Sites**
- Sites in configuration should NOT be charged
- Veolia's DCS team and customer success teams can configure without being charged
- Only production sites count toward licensing

#### 3. **User Management**
- Limit number of Veolia configuration users (currently ~12)
- Track HWF admin users
- Track by role: HWF admin, enterprise admin, plant users

#### 4. **Call Home Mechanism**
- Quarterly reporting to A-Stack license server
- Report active site counts, user counts
- Automated stats push from HWF to license server

---

## 🏗️ System Architecture

### Generic Component Design
- **License Server:** Must be built as a **generic component**, not HWF-specific
- **First Consumer:** HWF will be the first user
- **Future Consumers:** Boost and other on-prem applications

### License Server Requirements

#### 7 Core API Calls Required

1. **Create Site** - Create new site license
2. **Update Site** - Update site status/information
3. **Delete Site** - Remove site from licensing
4. **Refresh Key** - Monthly key refresh for security
5. **Get Aggregate Stats** - Quarterly reporting
6. **Check Validity** - Verify if site key is valid
7. **Send Alerts** - Notify when keys are invalid/expired

### Site Key Types

#### Dev vs Production Keys

**Production Keys:**
- Automatically assigned to all sites created on HWF side
- Tracked for billing/usage
- Includes: Basic plants, commissioning plants, active plants

**Dev Keys:**
- Only for sites originating from **Boost** with dev mode flag
- Boost explicitly requests dev key when creating site
- Not charged

#### Key Expiration & Refresh
- **Frequency:** Monthly refresh required (security requirement from Alex)
- **Process:** Client triggers refresh call → License server issues new key
- **Token Validity:** 1 month expiration
- **Enforcement:** If key is invalid/expired, site cannot be accessed

---

## 🔐 Key Management

### Key Hierarchical Structure

**Chain of Trust:**
```
Root (A-Stack) 
  ↓ Signs
Org (Hub)
  ↓ Signs  
Enterprise (Veolia)
  ↓ Signs
Site
```

**Enterprise-Level Keys Required:**
- Each enterprise needs its own API key
- Even if no production sites exist (only proxy sites)
- Prevents having enterprise with zero sites from going unnoticed

### Site Key Generation Flow

**For HWF Sites (Production):**
- Any site created on HWF side gets **production key** automatically
- No distinction between basic/commissioning/active
- All treated as production sites for licensing

**For Boost Sites:**
- **Dev Mode:** Boost specifies "dev" flag → gets dev key
- **Production Mode:** Boost specifies "production" flag → gets production key
- Key type determined at site creation time
- Can transition from dev to production via update API

### Key Storage & Security

**HWF Side:**
- Must use AWS Secrets Manager (requirement)
- Cannot store keys in database
- Keys cached with expiration

**License Server Side:**
- Can use database storage
- Or use AWS Key Management Service (optional)
- Up to implementation team

---

## 🔄 Workflows & Operations

### Site Provisioning Flow

1. **HWF creates site** → Makes API call to license server
2. **License server** creates site record with key
3. **HWF receives** API key
4. **Site uses key** to access HWF services

### Key Refresh Flow

1. **Client (Boost/HWF)** triggers refresh before expiration
2. **License server** validates old key
3. **License server** generates new key
4. **Client** receives new key
5. **Old key** becomes invalid

### Quarterly Stats Flow

1. **HWF** aggregates statistics (sites, users)
2. **HWF** pushes JSON file to license server
3. **A-Stack** retrieves stats from server
4. **Billing** calculated based on production sites

### Key Validation & Enforcement

**Validation Method:**
- Client-side caching with expiration
- Periodic callback to license server
- Token-based approach (cache for period, then refresh)

**When Key Invalid:**
- Site becomes inaccessible
- No data refresh from Insight
- Alert sent to A-Stack
- Same behavior as current Boost-to-HWF invalid key scenario

---

## 📊 Stats & Reporting

### Data Included in Manifest

**Site Data:**
- Total production sites
- Total dev sites (configuration only)
- Total enterprises
- Sites per enterprise breakdown

**User Data:**
- **HWF Admin users:** Names and email IDs
- **Enterprise Admin:** Total count only (no names)
- **Plant Users:** Total count only
- **Demo Users:** Total count only

### Privacy Considerations

**NOT Included:**
- Enterprise names
- Site names
- Individual user details (except HWF admins)
- Customer identifiable information

**Included:**
- Aggregate counts by category
- HWF admin contact information (names + emails)
- Enterprise-to-site mappings (counts)

### Counting Rules

**Included:**
- Commissioning sites ✅
- Basic sites ✅
- Active sites ✅
- Enterprise-level keys ✅

**NOT Included:**
- Proxy sites ❌
- Configuration-only sites (dev mode) ❌
- Nested duplicate sites (treated as separate) ✅

**Employee/Internal Users:**
- HWF admins (including A-Stack employees) counted
- Will transition to Veolia when they take over

---

## 🚀 Implementation Details

### API Modifications for Boost

**New API:** `Update API Key`
- Allows transition from dev to production
- Also handles monthly key refresh
- Same endpoint for both scenarios

**Modified API:** `Create Boost Dashboard`
- New parameter: `mode` (dev/production)
- Returns appropriate key type
- Returns site-specific API key

**Application-Level Key:**
- Boost needs application-level API key
- Separate from site-specific key
- Validates Boost as legitimate entity
- Call home to verify validity

### Update Requirements

**Nancy's Notes:**
- Update Boost dashboard API to include site key parameter
- Implement key refresh mechanism
- Support enterprise-level keys in chain of trust
- Handle commissioning sites as unique entities
- Implement stats aggregation and reporting

**Enhanced Requirements:**
- Quarterly reporting via JSON upload
- API key validation and refresh
- Dev/production key distinction
- Enterprise-level key management
- Invalid key alerting

---

## ⚠️ Important Clarifications

### Commissioning Sites Issue

**Problem:** Multiple sites can share same site ID but different states
- Example: Site ID "ABC123" with both commissioning and active states

**Solution:** Treat as separate sites
- Each state gets unique API key
- Both counted for licensing
- Customer's problem if they duplicate unnecessarily

### Configuration vs Production

**Challenge:** Hard to determine transition point
- No clear "configuration" vs "production" state in current system

**Solution:** Don't track state transitions
- All HWF sites = production
- Only Boost sites can be dev (when explicitly flagged)
- Simplified licensing model

### Watermark for Dev Keys

**Requirement:** Dev sites display watermark
- Indicates non-production environment
- Visual indication for configuration/testing
- Not yet implemented (future requirement)

---

## 📅 Next Steps & Timeline

### Immediate Actions

1. **Nancy:** Update requirements document with today's discussion
2. **An:** Create design document for implementation
3. **Kartik:** Review requirements and design (end of day Wednesday)
4. **Alok:** Approve design and get timeline estimate

### ETA Delivery

**Target:** Wednesday end of day
**Deliverables:**
- Updated requirements document
- Design document
- Implementation timeline
- Effort estimate

---

## 🎯 Key Decisions

### 1. Generic License Server Component
✅ Build as reusable component, not HWF-specific

### 2. Automated Provisioning
✅ No manual intervention required for site creation

### 3. Key Refresh Frequency
✅ Monthly refresh mandatory (security requirement)

### 4. Validation Approach
✅ Token-based with caching, periodic call home

### 5. Enforcement Behavior
✅ Same as current Boost-to-HWF: reject if invalid

### 6. Stats Reporting
✅ Quarterly JSON upload from HWF to license server

### 7. Dev vs Production
✅ Simple model: HWF = production, Boost = configurable

### 8. Enterprise Keys
✅ Each enterprise gets own key in chain of trust

### 9. Privacy
✅ Aggregate counts only, HWF admin contact info only

### 10. Commissioning Sites
✅ Treat as separate entities, count separately

---

## 🤔 Open Questions (Resolved)

✅ **Who enforces keys?** → HWF side only  
✅ **How often refreshed?** → Monthly  
✅ **Dev vs production?** → HWF = production, Boost = configurable  
✅ **Enterprise keys?** → Yes, required  
✅ **Commissioning sites?** → Treated separately  
✅ **Internal users?** → Counted as HWF admins  
✅ **Privacy?** → Aggregates only, no identifiable info  
✅ **UI vs API?** → Backend API only for now  
✅ **Watermarks?** → Future requirement  

---

## 📝 Summary

This transcript outlines a comprehensive licensing system for tracking and billing site usage across HWF and Boost platforms. The solution includes:

- Generic, reusable license server component
- Automated site provisioning and key management
- Monthly key refresh for security
- Quarterly stats reporting
- Dev/production key distinction
- Enterprise and site-level keys in chain of trust
- Privacy-compliant reporting
- Automated enforcement of valid keys

All participants agreed on the approach, with design document and timeline expected by end of day Wednesday.

