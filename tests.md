# Test Suite Specification: Workforce, Attendance & Payroll Platform

This document defines unit, integration, and end-to-end test cases mapped to the user stories in `task.md`. Each test case includes setup preconditions, execution steps, and expected assertions for an AI coding agent to implement automated tests.

---

## Epic 1: Multi-Tenant Foundation & RBAC

### US-01: Organization & Branch Onboarding
* [ ] **TC-01.1: Tenant Creation & Role Seeding**
  * **Type:** Integration
  * **Preconditions:** Empty database state.
  * **Steps:** Send `POST /api/v1/tenants` with valid company name, timezone (`Asia/Kolkata`), currency (`INR`), and logo file.
  * **Expected Result:** Returns `201 Created` with unique `tenant_id`; `Admin`, `Manager`, and `Employee` roles are automatically seeded for the tenant; logo URL is persisted.
* [ ] **TC-01.2: Branch Geofence & Wi-Fi Configuration**
  * **Type:** Integration
  * **Preconditions:** Authenticated `Admin` user for `Tenant A`.
  * **Steps:** Send `POST /api/v1/branches` with `lat: 28.6139`, `lng: 77.2090`, `radius_meters: 150`, and `wifi_bssids: ["00:11:22:33:44:55"]`.
  * **Expected Result:** Returns `201 Created`; branch coordinates, radius, and Wi-Fi BSSID list are stored accurately.
* [ ] **TC-01.3: Invalid Geofence Validation**
  * **Type:** Unit / API Validation
  * **Preconditions:** Authenticated `Admin` user.
  * **Steps:** Send `POST /api/v1/branches` with out-of-bounds `lat: 95.0` and `radius_meters: -10`.
  * **Expected Result:** Returns `400 Bad Request` with descriptive field validation errors.

### US-02: Branch & Role Permissions (RBAC)
* [ ] **TC-02.1: Scoped Manager Access Enforcement**
  * **Type:** Integration
  * **Preconditions:** `Manager M1` assigned strictly to `Branch A` with `approve_leaves: true`, `edit_attendance: true`, and `view_salary: false`.
  * **Steps:**
    1. Call `GET /api/v1/attendance?branch_id=Branch_A` as `M1`.
    2. Call `GET /api/v1/attendance?branch_id=Branch_B` as `M1`.
    3. Call `GET /api/v1/payroll/runs` as `M1`.
  * **Expected Result:** Step 1 returns `200 OK`; Steps 2 and 3 return `403 Forbidden`.
* [ ] **TC-02.2: Strict Cross-Tenant Isolation**
  * **Type:** Security / Integration
  * **Preconditions:** `Tenant A` and `Tenant B` exist with separate employees and branches.
  * **Steps:** Attempt `GET /api/v1/employees/:tenant_b_emp_id` using a valid JWT issued for `Tenant A`.
  * **Expected Result:** Returns `403 Forbidden` or `404 Not Found` with zero data leakage.

---

## Epic 2: Employee Lifecycle & Pay Configuration

### US-03: Employee Onboarding, Bulk Import & Document Vault
* [ ] **TC-03.1: Single Employee Creation & Presigned KYC Access**
  * **Type:** Integration
  * **Preconditions:** Authenticated `Admin` in `Tenant A`.
  * **Steps:** Create employee via `POST /api/v1/employees` with PAN, Aadhaar, UAN, ESIC, and upload a KYC document.
  * **Expected Result:** Unique `employee_code` generated; statutory IDs masked/encrypted at rest; document download endpoint returns a short-lived presigned URL.
* [ ] **TC-03.2: Bulk CSV/Excel Import Partial Error Reporting**
  * **Type:** Integration
  * **Preconditions:** Valid CSV file with 10 employee rows where Row 4 has an invalid PAN format and Row 7 duplicates an existing `employee_code`.
  * **Steps:** Upload file to `POST /api/v1/employees/import`.
  * **Expected Result:** Valid rows are processed cleanly and response returns row-level error diagnostics for Rows 4 and 7.
* [ ] **TC-03.3: Duplicate Employee Code Conflict**
  * **Type:** Integration
  * **Preconditions:** Employee with `employee_code: "EMP-001"` exists in `Tenant A`.
  * **Steps:** Attempt to create another employee with `employee_code: "EMP-001"` in `Tenant A`.
  * **Expected Result:** Returns `409 Conflict`.

### US-04: Flexible Pay-Basis & Salary Structure
* [ ] **TC-04.1: Multi-Basis Pay Configuration**
  * **Type:** Integration
  * **Preconditions:** Three active employees in `Tenant A`.
  * **Steps:** Assign `Monthly` (`₹45,000`) to Emp 1, `Daily Wage` (`₹850/day`) to Emp 2, and `Piece/Work-Rate` (`₹30/unit`) to Emp 3 via `POST /api/v1/employees/:id/salary-structure`.
  * **Expected Result:** All three pay structures persist with distinct `pay_basis` types and effective dates.
* [ ] **TC-04.2: Formula-Based Percentage Components**
  * **Type:** Unit
  * **Preconditions:** Salary structure with `Basic = ₹30,000` and `HRA = 40% of Basic`.
  * **Steps:** Evaluate salary structure preview.
  * **Expected Result:** `HRA` resolves deterministically to `₹12,000`.

### US-05: HR Letters & Full-and-Final (F&F) Settlement
* [ ] **TC-05.1: Dynamic Offer Letter PDF Generation**
  * **Type:** Integration
  * **Preconditions:** Active template containing `{{employee_name}}`, `{{designation}}`, and `{{joining_date}}`.
  * **Steps:** Call `POST /api/v1/employees/:id/letters` for target employee.
  * **Expected Result:** Returns `200 OK` with valid PDF binary stream containing substituted employee values and zero raw `{{...}}` tokens.
* [ ] **TC-05.2: Exit State Machine & Post-Exit Punch Block**
  * **Type:** Integration
  * **Preconditions:** Employee transitioned `Active` -> `Notice` -> `Exited` with `last_working_day: 2026-11-20`.
  * **Steps:** Attempt `POST /api/v1/attendance/punch` on `2026-11-21` for the exited employee.
  * **Expected Result:** Punch rejected with `403 Forbidden` / inactive employee business error.
* [ ] **TC-05.3: F&F Settlement Netting**
  * **Type:** Unit / Integration
  * **Preconditions:** Exiting employee has `5` encashable leave days (`+₹5,000`) and `₹2,000` pending advance balance.
  * **Steps:** Execute F&F settlement via `POST /api/v1/employees/:id/exit`.
  * **Expected Result:** Net F&F payable computes to `₹3,000` and pending advance balance updates to `₹0`.

---

## Epic 3: Shift Scheduling & Attendance Rules Engine

### US-06: Shift, Roster & Holiday Definition
* [ ] **TC-06.1: Cross-Midnight (Overnight) Shift Evaluation**
  * **Type:** Unit
  * **Preconditions:** Shift defined as `22:00` to `06:00` (next day).
  * **Steps:** Evaluate punch-in at `2026-11-23T22:04:00` and punch-out at `2026-11-24T06:02:00`.
  * **Expected Result:** Both punches group under the `2026-11-23` shift date with total worked duration of `7h 58m`.
* [ ] **TC-06.2: Roster Override Precedence**
  * **Type:** Unit
  * **Preconditions:** Branch default shift `09:30–18:30`; employee roster override `11:00–20:00` for `Nov 10–14`.
  * **Steps:** Resolve assigned shift for `Nov 12` and `Nov 15`.
  * **Expected Result:** `Nov 12` resolves to `11:00–20:00`; `Nov 15` reverts to `09:30–18:30`.
* [ ] **TC-06.3: Holiday & Alternate Week-Off Classification**
  * **Type:** Unit
  * **Preconditions:** Policy configured with 2nd & 4th Saturdays off plus a declared holiday on `Nov 14`.
  * **Steps:** Run daily attendance reconciliation on `Nov 14` (holiday) and 2nd Saturday with zero punches.
  * **Expected Result:** Dates evaluate to `Holiday` and `Weekly-Off` respectively (not `Absent`).

### US-07: Late, Early-Exit & Half-Day Policy Rules
* [ ] **TC-07.1: Grace Period vs. Late-Entry Fine**
  * **Type:** Unit
  * **Preconditions:** Shift starts at `09:30`, grace period `15 mins`, late fine `₹100`.
  * **Steps:**
    1. Evaluate punch-in at `09:42`.
    2. Evaluate punch-in at `09:50`.
  * **Expected Result:** Step 1 evaluates as `Present` with `₹0` fine; Step 2 evaluates with `Late` flag and `₹100` penalty item.
* [ ] **TC-07.2: Working-Hour Thresholds (Absent / Half-Day / Present)**
  * **Type:** Unit
  * **Preconditions:** Policy requires `>= 4.0 hrs` for `Half-Day` and `>= 8.0 hrs` for `Present`.
  * **Steps:** Evaluate total worked durations of `3.5 hrs`, `5.0 hrs`, and `8.2 hrs`.
  * **Expected Result:** `3.5 hrs` -> `Absent`; `5.0 hrs` -> `Half-Day`; `8.2 hrs` -> `Present`.

### US-08: Overtime & Comp-Off Rules
* [ ] **TC-08.1: Overtime Multiplier Calculation**
  * **Type:** Unit
  * **Preconditions:** Overtime threshold `30 mins` post-shift, multiplier `1.5x`.
  * **Steps:** Employee works `75 mins` past shift end and manager approves overtime.
  * **Expected Result:** `75 mins` recorded as approved overtime at `1.5x` hourly rate.
* [ ] **TC-08.2: Holiday Work to Comp-Off Conversion**
  * **Type:** Integration
  * **Preconditions:** Comp-off accrual enabled for working on holidays.
  * **Steps:** Employee punches full shift on a declared `Holiday`.
  * **Expected Result:** `1.0` `Comp-Off` credit added to employee's `leave_balances`.

---

## Epic 4: Multi-Mode Attendance Capture & Sync

### US-09: Mobile Selfie, Geofence & Wi-Fi Punching
* [ ] **TC-09.1: Valid Geofence & Selfie Punch**
  * **Type:** Integration
  * **Preconditions:** Branch geofence radius `150m`.
  * **Steps:** Submit `POST /api/v1/attendance/punch` from coordinates `45m` away with selfie URL.
  * **Expected Result:** Returns `201 Created`; punch recorded with `is_out_of_geofence: false`.
* [ ] **TC-09.2: Out-of-Geofence / Unapproved Wi-Fi Handling**
  * **Type:** Integration
  * **Preconditions:** Branch configured to flag out-of-radius punches.
  * **Steps:** Submit punch from `2.5 km` outside geofence and unrecognized Wi-Fi BSSID.
  * **Expected Result:** Punch flagged with `OUT_OF_GEOFENCE` for manager review (or rejected with `403` when strict-block mode is active).

### US-10: Kiosk Facial Recognition ("Lens" Mode)
* [ ] **TC-10.1: Face Enrollment & Verification**
  * **Type:** Integration
  * **Preconditions:** Employee enrolled via `POST /api/v1/kiosk/enroll`.
  * **Steps:** Submit matching face embedding to `POST /api/v1/kiosk/punch`.
  * **Expected Result:** Returns `201 Created` with matched `employee_id` and punch recorded.
* [ ] **TC-10.2: Duplicate Kiosk Punch Debounce**
  * **Type:** Integration
  * **Preconditions:** Employee punched in via kiosk `20 seconds` ago (`debounce_window = 120s`).
  * **Steps:** Submit second kiosk punch for the same employee immediately.
  * **Expected Result:** Second punch ignored/debounced; only one punch log persisted.

### US-11: Hardware Biometric Ingestion
* [ ] **TC-11.1: Idempotent Device Log Ingestion**
  * **Type:** Integration
  * **Preconditions:** Biometric device `BIO-01` mapped to `Emp A` (`device_user_id: "101"`).
  * **Steps:** Send identical payload `(device_id: "BIO-01", device_user_id: "101", punch_timestamp: "2026-11-23T09:00:00Z")` twice to `POST /api/v1/integrations/biometric/logs`.
  * **Expected Result:** Both requests succeed (`200/201`), but only one attendance record is stored.
* [ ] **TC-11.2: Unmapped Biometric User Quarantine**
  * **Type:** Integration
  * **Preconditions:** `device_user_id: "999"` has no employee mapping.
  * **Steps:** Ingest punch log for `"999"`.
  * **Expected Result:** Record stored in `unmapped_biometric_logs` without failing the batch ingestion request.

### US-12: Bulk Attendance, Work Output & Offline Sync
* [ ] **TC-12.1: Out-of-Order Offline Batch Replay**
  * **Type:** Integration
  * **Preconditions:** Offline mobile queue with `IN` (`09:00`, `uuid-1`) and `OUT` (`18:00`, `uuid-2`).
  * **Steps:** Send `uuid-2` first and `uuid-1` second (plus a duplicate `uuid-1`) to `POST /api/v1/attendance/sync-batch`.
  * **Expected Result:** Deduplicated by `client_event_uuid` and daily summary evaluates `09:00` as punch-in and `18:00` as punch-out.
* [ ] **TC-12.2: Bulk Override & Immutable Audit Trail**
  * **Type:** Integration
  * **Preconditions:** 10 employees marked `Absent` for `2026-11-22`.
  * **Steps:** Manager calls `POST /api/v1/attendance/bulk-override` changing all 10 to `Present` with reason `"Site power outage"`.
  * **Expected Result:** All 10 daily statuses update to `Present` and 10 entries are written to `attendance_audit_logs`.

---

## Epic 5: Leave, On-Duty, Advances & Expense Management

### US-13: Leave Policies, Multi-Level Approvals & On-Duty Requests
* [ ] **TC-13.1: Insufficient Leave Balance Rejection**
  * **Type:** Integration
  * **Preconditions:** Employee has `2.0` Casual Leave balance; LOP fallback disabled.
  * **Steps:** Submit leave request for `4.0` days of Casual Leave.
  * **Expected Result:** Returns `400 Bad Request` indicating insufficient leave balance.
* [ ] **TC-13.2: Approved Leave & On-Duty Attendance Lock**
  * **Type:** Integration
  * **Preconditions:** Pending `Paid Leave` for `Nov 24` and pending `On-Duty` request for `Nov 25`.
  * **Steps:** Manager approves both requests.
  * **Expected Result:** `Nov 24` locks as `Paid Leave` and balance decrements by `1.0`; `Nov 25` evaluates as `Present` without geofence flags.

### US-14: Mid-Month Advances & Loan EMIs
* [ ] **TC-14.1: Loan EMI Schedule Generation & Deduction**
  * **Type:** Integration
  * **Preconditions:** Active employee.
  * **Steps:** Approve `₹12,000` loan with `4` monthly installments of `₹3,000`.
  * **Expected Result:** 4 scheduled `loan_repayment_schedules` rows created; running 1 payroll cycle deducts `₹3,000` and updates remaining balance to `₹9,000`.

### US-15: Expense Claims & Site CashBook
* [ ] **TC-15.1: Site CashBook Running Balance**
  * **Type:** Integration
  * **Preconditions:** Branch CashBook starting at `₹0`.
  * **Steps:** Post `CASH_IN: ₹10,000`, `CASH_OUT: ₹2,500`, and `CASH_IN: ₹5,000`.
  * **Expected Result:** Current branch CashBook balance equals `₹12,500`.
* [ ] **TC-15.2: Approved Reimbursement Integration into Payroll**
  * **Type:** Integration
  * **Preconditions:** Employee submits `₹1,200` travel expense claim with receipt URL.
  * **Steps:** Manager approves claim; admin generates draft payroll for the cycle.
  * **Expected Result:** `₹1,200` appears as a non-taxable reimbursement line item in employee's draft payslip.

---

## Epic 6: Field Force GPS, Tasks & CRM Tracking

### US-16: Live Route Tracking & Anti-Tampering Alerts
* [ ] **TC-16.1: Route Breadcrumb & Stationary Stop Detection**
  * **Type:** Unit / Integration
  * **Preconditions:** Employee punched in for active shift.
  * **Steps:** Ingest batched pings via `POST /api/v1/tracking/pings` showing 30 mins of movement followed by 25 mins within a `30m` radius.
  * **Expected Result:** Route summary returns computed distance and identifies a `25 min` stationary stop.
* [ ] **TC-16.2: Mock Location & Tracking Gap Alerting**
  * **Type:** Integration
  * **Preconditions:** Active field employee shift.
  * **Steps:** Ingest ping with `is_mocked: true` and another ping after a `45 min` GPS-disabled gap.
  * **Expected Result:** Two `tracking_alerts` (`MOCK_LOCATION`, `GPS_DISABLED_GAP`) are created for manager visibility.

### US-17: Field Tasks, Client Visits & Geofenced Orders
* [ ] **TC-17.1: OTP-Verified Client Task Completion**
  * **Type:** Integration
  * **Preconditions:** Field task assigned at client location with OTP verification enabled.
  * **Steps:**
    1. Attempt task completion with invalid OTP.
    2. Submit valid client OTP.
  * **Expected Result:** Step 1 returns `400 Bad Request`; Step 2 marks task `Completed` with verified timestamp and coordinates.
* [ ] **TC-17.2: Client Visit Order Capture**
  * **Type:** Integration
  * **Preconditions:** Active client check-in session.
  * **Steps:** Submit field order with 2 line items (`3 x ₹500` and `2 x ₹250`).
  * **Expected Result:** Order persisted with total `₹2,000` linked to `visit_id`, `client_id`, and `employee_id`.

---

## Epic 7: Payroll Calculation, Statutory Compliance & Reporting

### US-18: Automated Payroll Run & Lock
* [ ] **TC-18.1: Full Payroll Calculation Pipeline**
  * **Type:** Integration
  * **Preconditions:** Monthly employee (`₹30,000` in 30-day month) with `2` LOP days (`-₹2,000`), `₹1,500` approved overtime, `₹200` late fine, `₹2,000` loan EMI, and `₹500` reimbursement.
  * **Steps:** Execute `POST /api/v1/payroll/runs` for the month.
  * **Expected Result:** Draft payslip calculates `Gross Earned = ₹29,500`, `Deductions = ₹2,200`, `Reimbursements = ₹500`, and `Net Payable = ₹27,800` with itemized `payslip_line_items`.
* [ ] **TC-18.2: Locked Payroll Immutability**
  * **Type:** Integration
  * **Preconditions:** Payroll run transitioned to `Locked` via `POST /api/v1/payroll/runs/:id/lock`.
  * **Steps:** Attempt to recalculate payroll or override attendance within the locked date range.
  * **Expected Result:** Rejected with `409 Conflict` / `403 Forbidden`.

### US-19: Indian Statutory Compliance & IT Declarations
* [ ] **TC-19.1: EPF, ESIC & State Professional Tax Calculation**
  * **Type:** Unit / Integration
  * **Preconditions:** Employee in Maharashtra (`PT` slab applicable) with `Basic = ₹15,000` (PF eligible) and `Gross = ₹18,000` (ESIC eligible).
  * **Steps:** Run payroll statutory engine.
  * **Expected Result:** Employee PF (`12%` = `₹1,800`), Employer PF (`12%` = `₹1,800`), Employee ESIC (`0.75%` = `₹135`), Employer ESIC (`3.25%` = `₹585`), and Maharashtra PT slab (`₹200`) compute accurately.
* [ ] **TC-19.2: IT Declaration Proof Verification & TDS Adjustment**
  * **Type:** Integration
  * **Preconditions:** Employee under Old Tax Regime submits `80C` declaration (`₹1,50,000`) and admin marks proofs `Verified`.
  * **Steps:** Run payroll TDS projection and export Form 24Q / ECR compliance reports.
  * **Expected Result:** Monthly TDS reflects verified `80C` deduction; `GET /api/v1/compliance/exports/ecr` and `form24q` return valid formatted outputs.

### US-20: Branded Payslips, Statutory Registers & Alerts
* [ ] **TC-20.1: Event-Driven Async Notifications**
  * **Type:** Integration
  * **Preconditions:** Mocked WhatsApp/SMS/Email adapters registered in queue worker.
  * **Steps:** Trigger `PunchRecorded` and `LeaveUpdated` domain events.
  * **Expected Result:** Notification jobs are enqueued asynchronously with expected recipient and template variables without blocking API response latency.
* [ ] **TC-20.2: Branded Payslip PDF & Statutory Register Exports**
  * **Type:** Integration
  * **Preconditions:** Locked payroll run with multiple employees.
  * **Steps:** Request `GET /api/v1/payslips/:id/pdf` and `GET /api/v1/reports/registers/muster-roll`.
  * **Expected Result:** Payslip PDF renders company logo, earnings/deductions, and YTD totals; Muster Roll export matches exact daily attendance counts and net pay figures.