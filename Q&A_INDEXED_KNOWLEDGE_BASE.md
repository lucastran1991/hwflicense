# Q&A Key Management - Indexed Knowledge Base

**Meeting Date:** October 27, 2025, 3:45 PM  
**Duration:** ~48 minutes  
**Participants:** Alok Batra, Kartik Shah, Nancy Tran, An Nguyen Thanh

---

## 📊 EXECUTIVE SUMMARY

A comprehensive licensing system to track and bill site usage across HWF (Hive Water Footprint) and Boost platforms for Veolia customer environment. The system enables A-Stack to monitor licensing without direct environment access.

---

## 🎯 BUSINESS REQUIREMENTS

### 1. **Problem Statement**
- Veolia will take over environment management within ~1 year
- A-Stack will lose direct access to their infrastructure
- Need to track active sites for billing purposes
- Must determine licensing costs based on:
  - Number of HWF sites
  - Number of Boost sites

### 2. **Key Business Goals**
- Track active production sites without environment access
- Avoid charging for configuration/testing sites
- Restrict configuration access to authorized Veolia teams
- Quarterly licensing calculation (not continuous)
- Limit number of Veolia configuration users (~12 people)
- Monitor HWF admin users

### 3. **Challenges**
- Hard to determine when site transitions from "configuration" to "production"
- Veolia provides full service - customers may never access software directly
- Sites can exist in semi-configured state but still be "active"

---

## 🏗️ SYSTEM ARCHITECTURE

### Core Principle: Generic Reusable Component
**Decision:** License server must be built as a **generic component**, not HWF-specific

### First Consumer
- HWF (Hive Water Footprint) - first implementation
- Future: Boost and other on-prem applications

### License Server APIs (7 Core Endpoints)

| # | API | Purpose |
|---|-----|---------|
| 1 | **Create Site** | Create new site license |
| 2 | **Update Site** | Update site status/information |
| 3 | **Delete Site** | Remove site from licensing |
| 4 | **Refresh Key** | Monthly key refresh for security |
| 5 | **Get Aggregate Stats** | Quarterly reporting |
| 6 | **Check Validity** | Verify site key is valid |
| 7 | **Send Alerts** | Notify when keys are invalid/expired |

---

## 🔐 KEY MANAGEMENT

### Key Types

#### Production Keys
- **Assigned to:** All sites created on HWF side (automatically)
- **Includes:** Basic plants, commissioning plants, active plants
- **Billed:** Yes, counted toward licensing

#### Dev Keys
- **Assigned to:** Only sites originating from **Boost** with dev mode flag
- **Criteria:** Boost must explicitly request "dev" when creating site
- **Not billed:** Configuration/testing sites

### Key Lifecycle

**Creation Flow:**
1. HWF creates site → Calls license server
2. License server generates appropriate key (dev/prod)
3. HWF receives and stores key
4. Site uses key for all operations

**Refresh Flow:**
- **Frequency:** Monthly (security requirement from Alex)
- **Process:** Client triggers refresh → License server validates → Issues new key
- **Expiration:** Keys expire after 1 month
- **Enforcement:** Invalid keys block site access

**Hierarchical Chain:**
```
Root (A-Stack) 
  ↓ Signs
Org (Hub) 
  ↓ Signs  
Enterprise (Veolia)
  ↓ Signs
Site
```

### Enterprise Keys
**Decision:** Each enterprise requires its own API key
- Even if no production sites exist
- Prevents enterprises with only proxy sites from being missed
- Added to chain of trust hierarchy

---

## 🌐 SITE PROVISIONING

### HWF Sites (All Production)
- Any site created on HWF = **automatic production key**
- No dev option for HWF-originating sites
- Includes: basic, commissioning, active plants
- No distinction between plant states

### Boost Sites (Configurable)
- Boost can request "dev" or "production"
- Dev keys = configuration/testing (not billed)
- Production keys = active sites (billed)
- Can transition between dev → production via Update API

### Site States Included
- ✅ Commissioning sites
- ✅ Basic sites  
- ✅ Active sites
- ❌ **NOT included:** Proxy sites

### Commissioning Sites Edge Case
**Problem:** Multiple sites can share same site ID but different states
**Solution:** Treat as separate entities
- Each state gets unique API key
- Both counted for licensing
- Customer's problem if they duplicate

---

## 📊 STATS & REPORTING

### Data Included in Quarterly Manifest

#### Site Statistics
- Total production sites
- Total dev sites
- Sites per enterprise breakdown
- Total enterprises

#### User Statistics
**HWF Admins (Full Info):**
- Names + email IDs
- Admin contact information

**Aggregates Only:**
- Enterprise admin count
- Plant users count
- Demo users count
- Other role counts

### Privacy Rules

**NOT Included:**
- Enterprise names
- Site names
- Individual user emails (except HWF admins)
- Customer identifiable information
- Organization-specific breakdowns

**Included:**
- Aggregate counts
- HWF admin contact info
- Enterprise-to-site mappings (numbers only)

### Special Cases

**Internal Users (A-Stack Employees):**
- Counted as HWF admins
- Names and emails included
- Will transition to Veolia when they take over (in ~1 year)

---

## ⚙️ VALIDATION & ENFORCEMENT

### Validation Method
- **Token-based with caching**
- **Periodic call home** to license server
- **Token expiration:** 1 month validity
- **Cache refresh:** Before token expires

### Invalid Key Behavior
When site key is invalid:
1. Site becomes inaccessible
2. No data refresh from Insight
3. API calls rejected (same as current Boost-to-HWF behavior)
4. Alert sent to A-Stack license server
5. A-Stack gets notification

**Result:** Site cannot be accessed until valid key provided

---

## 🔄 WORKFLOWS

### Workflow 1: New Site Creation

```
1. HWF creates site in system
2. HWF calls License Server: Create Site API
3. License Server generates production key
4. HWF stores key in AWS Secrets Manager
5. Site uses key for all operations
```

### Workflow 2: Boost Site Creation (Dev)

```
1. Boost determines site needs "dev" mode
2. Boost calls: Create Boost Dashboard (with mode=dev)
3. HWF receives request → Calls License Server
4. License Server returns dev key
5. Boost stores key
6. Site uses dev key (watermark displayed)
```

### Workflow 3: Monthly Key Refresh

```
1. Client detects key approaching expiration
2. Client calls: Refresh Key API
3. License Server validates old key
4. License Server generates new key
5. Client receives new key
6. Old key becomes invalid
```

### Workflow 4: Quarterly Stats Report

```
1. HWF aggregates all statistics
2. HWF creates JSON manifest
3. HWF calls: Send Stats API
4. Stats stored on A-Stack license server
5. A-Stack retrieves stats for billing
```

### Workflow 5: Site Status Update

```
1. User marks site as production
2. Boost calls: Update API Key
3. HWF receives request
4. HWF calls License Server: Update Site
5. Key transitions dev → production
6. Site now counted for billing
```

---

## 🔒 SECURITY REQUIREMENTS

### Key Storage

**HWF Side (Mandatory):**
- **MUST use:** AWS Secrets Manager
- Cannot store keys in database
- Requirement from Kartik/Alok

**License Server Side (Optional):**
- Can use database storage
- Can use AWS Key Management Service
- Up to implementation team

### Key Properties

**Monthly Refresh:**
- Security requirement from Alex
- Keys must not be constant
- Monthly rotation mandatory

**Token Validity:**
- 1 month expiration
- Aligned with refresh frequency
- Client must refresh before expiration

---

## 📝 API SPECIFICATIONS

### Application-Level Keys (Boost)

**Requirement:** Two-tier key system
1. **Application Key:** Validates Boost as legitimate entity
2. **Site Key:** Site-specific key for operations

**Flow:**
- Boost uses application key to call HWF
- HWF validates application key
- Returns site-specific key
- Both keys managed by license server

**Rationale:** 
- Industry standard practice
- Prevents rogue clients
- Centrally managed keys

### Boost APIs to Modify

**1. Create Boost Dashboard**
- Add parameter: `mode` (dev/production)
- Return: appropriate key type
- Return: site-specific API key

**2. New: Update API Key**
- Update key: dev → production or vice versa
- Also handles monthly key refresh
- Same endpoint for both scenarios

**3. Existing APIs (No changes):**
- Get Asset Mapping External
- Update Asset Parameter by ID

---

## ⚠️ IMPORTANT DECISIONS

### 1. Generic vs HWF-Specific
✅ **Decision:** Build as generic reusable component

### 2. Dev vs Production Keys
✅ **Decision:** HWF = production, Boost = configurable

### 3. Key Refresh Frequency
✅ **Decision:** Monthly mandatory refresh

### 4. Validation Approach
✅ **Decision:** Token-based with caching

### 5. Enforcement Behavior
✅ **Decision:** Same as current Boost-to-HWF (reject if invalid)

### 6. Commissioning Sites
✅ **Decision:** Counted separately, treated as unique entities

### 7. Enterprise Keys
✅ **Decision:** Required in chain of trust

### 8. Privacy
✅ **Decision:** Aggregates only, HWF admin contact only

### 9. Internal Users
✅ **Decision:** Counted as HWF admins (will transition to Veolia)

### 10. Proxy Sites
✅ **Decision:** NOT counted for licensing

### 11. Watermarking
📋 **Future:** Dev keys should display watermark

---

## 🚀 IMPLEMENTATION DETAILS

### Team Composition
**Current Team (3 people):**
- Nancy Tran - Requirements & Design
- An Nguyen Thanh - Development
- Long - Currently 100% on other projects (not available yet)

### Deliverables Timeline

**Target Review Date:** End of day Wednesday

**Nancy's Tasks:**
1. Update requirements document with today's discussion
2. Send to Kartik for review

**An's Tasks:**
1. Create design document
2. Outline implementation approach
3. Provide effort estimate

**Kartik's Tasks:**
1. Review requirements document
2. Review design document
3. Get Alok's approval
4. Submit timelines

### Key Next Steps

1. ✅ Requirements clarified in meeting
2. 📝 Update requirements document (Nancy)
3. 📐 Create design document (An)
4. 👀 Review both documents (Kartik)
5. ✅ Get approval from Alok
6. 📊 Provide timeline/effort estimate
7. 🏗️ Begin implementation

---

## 📋 TECHNICAL REQUIREMENTS

### Database Considerations
- Store key metadata (expiration, type, site ID)
- Track site states
- Support quarterly stats retrieval
- Audit trail for key operations

### API Design
- RESTful endpoints
- Authentication required
- Return appropriate key type based on request
- Support batch operations where applicable

### Integration Points
- HWF → License Server (automated)
- License Server → A-Stack (quarterly stats)
- Boost → HWF → License Server (dev keys)

### Error Handling
- Invalid keys → block access
- Expired keys → force refresh
- Network failures → retry logic
- Alert on violations

---

## 💡 KEY TAKEAWAYS

### Business Perspective
- Enables accurate billing without environment access
- Automates license tracking
- Supports configuration vs production distinction
- Quarterly billing model

### Technical Perspective
- Generic reusable component
- Secure key management
- Automated provisioning
- Token-based validation
- Privacy-compliant reporting

### Operational Perspective
- Monthly key refresh (security)
- Quarterly stats reporting
- Automatic enforcement
- Alert on violations
- Simple key rotation

---

## ❓ UNRESOLVED QUESTIONS

None - All questions resolved during meeting.

Previously discussed but now clarified:
- ✅ Generic vs specific → Generic
- ✅ Dev vs production → Configurable for Boost
- ✅ Validation method → Token with caching
- ✅ Enforcement → Block access
- ✅ Privacy → Aggregates only
- ✅ Enterprise keys → Required
- ✅ Commissioning sites → Separate counting
- ✅ Internal users → Count as HWF admins

---

## 📞 MEETING PARTICIPANT ROLES

**Alok Batra:**
- Business requirements owner
- Final decision maker
- Authored initial requirements

**Kartik Shah:**
- Technical architect
- Solution design
- API specifications lead

**Nancy Tran:**
- Requirements coordinator
- Technical implementation lead
- Q&A clarifications

**An Nguyen Thanh:**
- Design and implementation
- Development lead
- Timeline estimation

---

## 🎓 LESSONS LEARNED

### Design Principles
1. Build for reuse, not single-use
2. Security first - monthly refresh mandatory
3. Privacy compliant - aggregates only
4. Automated workflows - no manual intervention
5. Token-based validation - efficient and secure

### Business Alignment
1. Track production sites accurately
2. Don't bill for configuration/testing
3. Support dev vs production distinction
4. Quarterly billing model
5. Future-proof for Veolia takeover

### Technical Excellence
1. Generic component architecture
2. Hierarchical key management
3. Secure key storage (AWS Secrets Manager)
4. Token-based validation with caching
5. Automated provisioning and refresh

---

## 📅 AGREED ACTIONS

### Tomorrow (Next Day)
- [ ] Nancy: Update requirements document
- [ ] Nancy: Send updated requirements to Kartik
- [ ] An: Create design document
- [ ] Nancy: Let Alok know ETA will be provided

### End of Day Wednesday
- [ ] Kartik: Review requirements document
- [ ] Kartik: Review design document
- [ ] Kartik: Get Alok's approval
- [ ] All: Submit timeline/effort estimate

### After Approval
- [ ] Begin implementation
- [ ] Follow design document
- [ ] Deliver to requirements

---

END OF KNOWLEDGE BASE
