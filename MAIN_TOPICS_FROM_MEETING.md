# Main Topics From Q&A Meeting

**Meeting Date:** October 27, 2025, 3:45 PM  
**Duration:** 48 minutes  
**Participants:** Alok Batra, Kartik Shah, Nancy Tran, An Nguyen Thanh

---

## 📋 MAIN TOPICS DISCUSSED

### 1. **Background & Business Context** ⏱️ 0:28 - 6:25
- Why licensing system is needed (Veolia takeover in ~1 year)
- Problem: A-Stack will lose access to environment
- Need to track site usage for billing (HWF + Boost sites)
- Requirement: Quarterly licensing calculation
- Goal: Avoid charging for configuration/testing sites
- Limit: ~12 Veolia configuration users
- Challenge: Hard to determine configuration vs production transition

### 2. **System Architecture** ⏱️ 6:25 - 8:22
- **Generic reusable component** (not HWF-specific)
- License server design - 7 core APIs needed
- Site key structure: HWF ↔ License Server communication
- Automated provisioning flow (no manual intervention)

### 3. **Key Types: Dev vs Production** ⏱️ 8:22 - 14:00
- **Production keys**: All HWF sites (automatic)
- **Dev keys**: Only from Boost when explicitly requested
- Plant states: basic, commissioning, active (all production for HWF)
- HWF admin vs Enterprise admin distinction
- User counting and privacy considerations

### 4. **Key Refresh & Security** ⏱️ 14:00 - 18:14
- Monthly key refresh required (Alex's requirement)
- Security rationale: Keys must not be constant
- Client-triggered refresh flow
- AWS Secrets Manager requirement for HWF side
- License server storage options (flexible)

### 5. **Enforcement & Validation** ⏱️ 18:14 - 22:57
- Who enforces: HWF side only
- Validation method: Token-based with caching
- Call home mechanism (periodic)
- Invalid key behavior: Block access, send alert
- Same behavior as current Boost-to-HWF

### 6. **Boost Integration** ⏱️ 22:57 - 32:20
- Two-tier API key system (application + site level)
- Create Boost Dashboard API modification
- Update API Key endpoint (dev ↔ production)
- Application-level key for Boost authentication
- Monthly refresh integration

### 7. **Enterprise & Site Management** ⏱️ 32:20 - 38:18
- Chain of trust: Root → Org → Enterprise → Site
- Enterprise-level keys required (even with proxy sites only)
- Commissioning sites counted separately
- Proxy sites NOT counted
- Duplicate site IDs handled

### 8. **Stats & Privacy** ⏱️ 38:18 - 42:31
- Quarterly manifest content
- HWF admin info: Names + emails (full disclosure)
- Enterprise admin, plant users: Counts only
- Enterprise names: NOT included (privacy)
- Internal users (A-Stack employees) counted as HWF admins
- Transition to Veolia in ~1 year

### 9. **Implementation Planning** ⏱️ 42:31 - 48:22
- Team: Nancy, An, Long (not available yet)
- Deliverables: Requirements doc + Design doc
- Timeline: Review by end of day Wednesday
- Tasks assigned: Nancy updates requirements, An creates design
- Alok approval needed before implementation
- ETA to be provided

---

## 🎯 KEY DECISIONS MADE

### Architecture
1. ✅ Generic reusable component (not HWF-specific)
2. ✅ Automated workflows (no manual intervention)
3. ✅ Token-based validation with caching
4. ✅ HWF side enforcement only

### Key Management
5. ✅ Monthly refresh mandatory (security)
6. ✅ Dev keys only from Boost (explicit request)
7. ✅ Enterprise-level keys required in chain
8. ✅ AWS Secrets Manager for HWF side

### Business Rules
9. ✅ Quarterly billing model
10. ✅ Configuration sites not billed (dev keys)
11. ✅ Proxy sites not counted
12. ✅ Commissioning sites counted separately

### Privacy & Reporting
13. ✅ Aggregate counts only (privacy)
14. ✅ HWF admin contact info included
15. ✅ Enterprise names excluded
16. ✅ Internal users counted as HWF admins

### Technical
17. ✅ 7 core API endpoints identified
18. ✅ Two-tier key system for Boost
19. ✅ Update API for dev ↔ production transition
20. ✅ Watermark for dev keys (future requirement)

---

## 💡 CRITICAL INSIGHTS

### Problem Understanding
- Veolia will manage environment independently in ~1 year
- No longer have access to track site usage directly
- Need automated license tracking and enforcement
- Quarterly billing based on production sites only

### Solution Approach
- Build once, use everywhere (generic component)
- Automatic provisioning and validation
- Monthly key refresh for security
- Quarterly stats reporting
- Dev vs production key distinction

### Business Alignment
- Track production sites accurately
- Don't charge for configuration/testing
- Support Veolia's full-service model
- Future-proof for environment takeover

### Technical Requirements
- 7 API endpoints minimum
- AWS Secrets Manager integration
- Token-based validation
- Automatic enforcement
- Privacy-compliant reporting

---

## 📊 TOPICS BY TIME

| Time | Topic | Duration |
|------|-------|----------|
| 0:00-0:28 | Meeting start | - |
| 0:28-6:25 | Background & requirements | 6 min |
| 6:25-8:22 | System architecture | 2 min |
| 8:22-10:29 | API requirements (create, update, delete, stats) | 2 min |
| 10:29-14:00 | Dev vs production keys | 4 min |
| 14:00-17:19 | Key refresh, security | 3 min |
| 17:19-22:57 | Enforcement, validation | 6 min |
| 22:57-26:02 | Boost integration | 3 min |
| 26:02-32:20 | Two-tier key system | 6 min |
| 32:20-38:18 | Enterprise keys, commissioning | 6 min |
| 38:18-42:31 | Stats, privacy, users | 4 min |
| 42:31-48:22 | Implementation planning | 6 min |

**Total Discussion:** ~48 minutes

---

## 🎓 SUMMARY BY PARTICIPANT CONTRIBUTION

### Alok Batra (Business Owner)
- Set business context and requirements
- Decided on dev vs production key logic
- Clarified enforcement behavior
- Approved approach

### Kartik Shah (Technical Architect)
- Specified system architecture
- Identified 7 API requirements
- Designed key management approach
- Owned enforcement strategy

### Nancy Tran (Implementation Lead)
- Clarified requirements details
- Asked about integration points
- Addressed Boost APIs
- Will update requirements document

### An Nguyen Thanh (Developer)
- Confirmed understanding
- Will create design document
- Timeline estimation

---

## ✅ RESOLUTION STATUS

### ✅ Resolved Topics
- Generic vs specific architecture
- Dev vs production key model
- Key refresh frequency
- Enforcement approach
- Validation method
- Enterprise key requirements
- Privacy rules
- Commissioning sites handling
- Internal users counting
- Proxy sites exclusion

### 📋 Action Items
- Nancy: Update requirements document by tomorrow
- An: Create design document by Wednesday
- Kartik: Review and approve by Wednesday
- All: Provide timeline/effort estimate

---

## 🚀 NEXT STEPS

1. **Requirements document** updated with all decisions
2. **Design document** created with technical specs
3. **Review meeting** with Kartik
4. **Alok approval** for approach
5. **Implementation** begins after approval
6. **Timeline** submitted to stakeholders

---

**Meeting Status:** ✅ COMPLETE - All topics discussed, all questions answered, all decisions made.

