# `SPEC.md` — Attendance & Payroll Super-Project (`attendance_management`)
**Architecture Style:** Go Workspaces (`go.work`) Modular Monolith (Single Binary `cmd/server` with Selective Module Activation via CLI/Env)  
**Root Go Module:** `attendance_management` (`attendance_management/pkg/...`, `attendance_management/modules/<name>`)  
**Taskwarrior Project:** `attendance_management1`  
**Backend Stack:** Go 1.23+, `gin-gonic/gin` (HTTP Server, RouterGroups & Middleware), `pgx/v5` (`pgxpool`), `pressly/goose/v3` (embedded SQL migrations), Headless Keycloak API (`golang.org/x/oauth2` + JWKS verification), `log/slog`, `caarlos0/env`  
**Testing:** `stretchr/testify`, `testcontainers-go/modules/postgres`, `uber-go/mock`, `net/http/httptest`  
**Web Frontend:** Angular + Angular Material UI  
**Mobile & Gate Kiosk Clients:** Flutter + Material 3 UI (Personal Employee App + Gate Kiosk FaceScan/FingerScan App)  

---

## Part 1: AI Agent Operational Guardrails (Zero-Hallucination & Scope Lock)

### 1. Grounding & Dependency Lock
- **Approved Stack Only:** NEVER introduce unapproved web frameworks, CLI libraries (use Go's standard `flag` package alongside `caarlos0/env`), or ORMs. Use strictly `gin-gonic/gin`, `pgx/v5`, `pressly/goose/v3`, `golang.org/x/oauth2`, `log/slog`, `caarlos0/env`, `stretchr/testify`, `testcontainers-go`, and `uber-go/mock`.
- **Verify Before Referencing:** NEVER import or reference a Go module, package, SQL table/column, Gin context key, Angular service, or Flutter widget unless verified in the workspace via file read or search during this session.
- **No Hand-Edited Applied Migrations:** Never modify already-committed SQL migration files; add sequential `00N_<description>.sql` files inside the owning module's `migrations/` directory with `-- +goose Up` and `-- +goose Down` annotations.

### 2. Strict Architectural & Anti-Diversion Guardrails
1. **One-Way Dependency Rule:**
   - `pkg/` (root SDK `attendance_management/pkg/...`) **MUST NEVER** import `modules/*`.
   - Sibling feature modules (`attendance_management/modules/organization`, `attendance_management/modules/attendance`, `attendance_management/modules/payroll`, `attendance_management/modules/user`) **MUST NEVER** import each other.
   - Only `cmd/server/main.go` imports both `pkg/` and `modules/*` via blank imports (`_`) so every module self-registers into `pkg/app`'s registry via `init()`.
2. **Single Binary Selective Activation (`cmd/server`):**
   - `cmd/server` compiles all modules into one binary and uses `-modules=organization,attendance` (or `ENABLED_MODULES` env var) to filter which modules run `Migrate()` and `Mount()`. If empty/omitted, **all** registered modules are activated.
3. **Gin Handler Isolation (`module.go` vs. `service.go`):**
   - `*gin.Context` **MUST NEVER** be passed into `service.go` or `repository.go`. Extract `context.Context` (`c.Request.Context()`), `TenantID`, and `UserClaims` in the Gin handler and pass pure Go types to the service layer.
4. **PostgreSQL Schema Isolation & Multi-Tenancy:**
   - Each module owns its dedicated PostgreSQL schema (`org.*`, `attendance.*`, `payroll.*`) and embeds its own `//go:embed migrations/*.sql`. Cross-schema SQL `JOIN`s and cross-schema foreign keys are **strictly forbidden**.
   - Every table MUST include `tenant_id TEXT NOT NULL`, and every SQL query (`SELECT`, `UPDATE`, `DELETE`) MUST filter by `tenant_id`.
5. **Decoupled Event-Driven Integration:**
   - Cross-module workflows MUST go through `core.Events.Publish(...)` / `Subscribe(...)`. When running a subset of modules via `-modules`, modules must gracefully handle missing peer subscribers without panicking.
6. **Minimal Diff & Plan Lock:**
   - Before writing code for any task, output the exact list of files to create/modify and confirm which module/layer they belong to. Do not refactor unrelated modules.

### 3. Mandatory Verification Loop
- **Backend:** Run `go work sync`, `go vet ./...`, regenerate mocks (`mockgen`) if interfaces changed, and execute `go test -race ./...` (including `httptest` Gin handler tests and `testcontainers-go` PostgreSQL suites).
- **Web Frontend:** Run `ng build` and verify Angular Material imports compile without errors.
- **Mobile App:** Run `flutter analyze` and `flutter test`.

---

## Part 2: High-Level System Architecture

### 1. Single-Binary Super-Project Topology (`go.work` + `cmd/server`)

```text
+-----------------------+   +-----------------------+   +-----------------------+
|   Angular Web App     |   |  Flutter Employee App |   |  Flutter Gate Kiosk   |
| (Material UI - Native |   | (Material 3 - Native  |   | (Material 3 - Face &  |
|  Login & Admin/Portal)|   |  Login, Punch & Logs) |   |  FingerScan Terminal) |
+-----------------------+   +-----------------------+   +-----------------------+
            |                           |                           |
            +---------------------------+---------------------------+
                                        | HTTPS / REST JSON (Bearer JWT)
                                        v
+-------------------------------------------------------------------------------+
|               SINGLE ENTRYPOINT: cmd/server/main.go                           |
|   CLI Flag: -modules=attendance,payroll   (Default/Empty => Run ALL Modules)  |
|   Blank-imports all modules -> init() registers -> app filters active list    |
+-------------------------------------------------------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
|                     PARENT SDK: pkg/ (attendance_management)                   |
|  pkg/app (Module Registry, Selective Filter, gin.Engine, /api/<module> Mount) |
|  pkg/httpx (Headless Auth Proxy, JWKS Verify, RequireRole, BindJSON[T], slog) |
|  pkg/events (Thread-safe InMemoryBus for cross-module Pub/Sub)                |
|  pkg/db (pgxpool, TenantTx Isolation Helper, Schema-Scoped goose Runner)      |
+-------------------------------------------------------------------------------+
         |                              |                              |
         v (Active if enabled)          v (Active if enabled)          v (Active if enabled)
+--------------------+         +--------------------+         +--------------------+
| modules/           |         | modules/           |         | modules/           |
|   organization     |         |   attendance       |         |   payroll          |
|  (Schema: org.*)   |         | (Schema:           |         | (Schema:           |
|                    |         |  attendance.*)     |         |  payroll.*)        |
+--------------------+         +--------------------+         +--------------------+
         |                              |                              |
         +------------------------------+------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
|        PostgreSQL 16 (Isolated Schemas)  |  Keycloak (Headless API Server)    |
+-------------------------------------------------------------------------------+
```

---

### 2. Gin Core Contracts (`pkg/app` & `pkg/httpx`)

* **`FeatureModule` Interface (`pkg/app/module.go`):**
  ```go
  type FeatureModule interface {
      Name() string
      Migrate(ctx context.Context, pool *pgxpool.Pool) error
      Mount(core *CoreContext, rg *gin.RouterGroup)
  }
  ```
* **Selective Module Activation (`pkg/app/module.go` & `cmd/server/main.go`):**
  * Parses `-modules` via Go's standard `flag.String("modules", "", "comma-separated list of modules to run; empty runs all")`, falling back to `ENABLED_MODULES` from `caarlos0/env`.
  * `app.ResolveModules(requested []string) ([]FeatureModule, error)`:
    * Empty slice returns all registered modules sorted deterministically by `Name()`.
    * Unknown module names fail fast at startup with an error listing valid registered modules.
  * `GET /healthz` returns `{"status":"ok","active_modules":["organization","attendance","payroll"]}`.
* **Typed Gin Helpers (`pkg/httpx/response.go`):**
  * `httpx.BindJSON[T any](c *gin.Context) (T, bool)` — runs `c.ShouldBindJSON(&req)` and aborts with standardized RFC 7807 `400 Bad Request` JSON on failure.
  * `httpx.JSON(c *gin.Context, status int, data any)` and `httpx.Error(c *gin.Context, status int, code, detail string)`.

---

### 3. Authentication & Roles (Headless Keycloak API Server)

* **No Browser Redirects:** Keycloak acts strictly as a headless OAuth2/OIDC token and Admin REST API server. Angular and Flutter render native Material login forms and **never** open webviews or redirect browsers to Keycloak's hosted login pages.
* **Backend Auth Proxy Endpoints (`pkg/httpx/auth_proxy.go` — unauthenticated `/api/auth/*`):**
  * `POST /api/auth/login` — Exchanges `{username, password}` with Keycloak's token endpoint (`/protocol/openid-connect/token` via `golang.org/x/oauth2` Direct Access Grant) using the server-side confidential client credentials and returns `{access_token, refresh_token, expires_in}`.
  * `POST /api/auth/refresh` — Exchanges `{refresh_token}` via Keycloak's token endpoint without exposing `KEYCLOAK_CLIENT_SECRET` to clients.
  * `POST /api/auth/logout` — Revokes the refresh token against Keycloak's `/protocol/openid-connect/revoke` endpoint.
* **Keycloak Admin REST API (`modules/organization`):**
  * When an `ORG_ADMIN` provisions an employee (`POST /api/organization/employees`), the service obtains a service-account token (`golang.org/x/oauth2/clientcredentials`), creates the user in Keycloak via `/admin/realms/{realm}/users`, assigns the `EMPLOYEE` or `GATE_KIOSK` role, and saves the returned Keycloak `sub` UUID in `org.employees`.
* **Local JWKS Verification (`pkg/httpx/middleware.go`):**
  * All `/api/<module>/*` routes are guarded by `httpx.KeycloakAuth(...)`, which caches RSA public keys from `{KEYCLOAK_ISSUER}/protocol/openid-connect/certs`, verifies RS256 signature/expiry/issuer locally, injects `sub`, `org_id` (`TenantID`), and roles into `gin.Context`, and enforces RBAC via `httpx.RequireRole("ORG_ADMIN", "EMPLOYEE", "GATE_KIOSK")`.

---

### 4. Core Domain Specifications

#### A. Organization Attendance Policy (`modules/organization`, Schema: `org`)
* Stores per-tenant policy in `org.policies`:
  * Allowed punch methods (`TEXT[]`): `FACESCAN_KIOSK`, `FINGERSCAN_KIOSK`, `WEB_LOGIN`, `MOBILE_LOGIN`.
  * `salary_slip_enabled` (`BOOLEAN`) — enables or disables salary slip generation in `modules/payroll`.
  * Shift boundaries (`shift_start_time`, `shift_end_time`).
* Emits `"org.policy.updated"` onto `core.Events` whenever policy settings change.

#### B. Attendance Tracking & Reporting (`modules/attendance`, Schema: `attendance`)
* Subscribes to `"org.policy.updated"` to maintain a local `attendance.tenant_policies` table (zero cross-schema joins).
* **Endpoints (`/api/attendance/*`):**
  * `POST /punch` (`EMPLOYEE` or `GATE_KIOSK`): Validates that `punch_method` is enabled for the tenant, rejects duplicate consecutive `IN` without `OUT` (and `OUT` without an open `IN`), records raw event in `attendance.logs`, pairs `IN`/`OUT` in `attendance.paired_sessions` with calculated `duration_minutes`, and publishes `"attendance.session.completed"`.
  * `GET /me` (`EMPLOYEE`): Lists the authenticated user's punch logs and paired sessions (filterable by date range).
  * `GET /reports` (`ORG_ADMIN`): Aggregated organization-wide attendance report (present days, total worked minutes, late/early flags) filterable by date range and employee.

#### D. User Management (`modules/user`, Schema: `user`)
* Stores user personal details and preferences
* **Endpoints (`/api/user/*`):**
  * `GET /me` - Get current user details
  * `PUT /me` - Update current user details
  * `GET /me/preferences` - Get user preferences
  * `PUT /me/preferences` - Update user preferences

---

### 5. Canonical Monorepo Directory Structure

```text
attendance_management/
├── go.work                                 # Links root SDK + all child feature modules
├── go.mod                                  # Parent SDK: "attendance_management"
├── docker-compose.yml                      # PostgreSQL 16 + Keycloak (auto-imported realm)
├── infra/
│   └── keycloak/realm-attendance.json      # Realm roles (ORG_ADMIN, EMPLOYEE, GATE_KIOSK) & client
│
├── cmd/
│   └── server/
│       └── main.go                         # Single entrypoint: parses -modules flag, wires core & active modules
│
├── pkg/                                    # Shared parent infrastructure (NO module imports)
│   ├── config/
│   │   └── env.go                          # caarlos0/env typed config loader (incl. ENABLED_MODULES)
│   ├── app/
│   │   ├── module.go                       # CoreContext, FeatureModule interface, Register(), ResolveModules()
│   │   └── server.go                       # gin.Engine setup, /api/auth/*, /api/<module> groups, /healthz
│   ├── db/
│   │   ├── postgres.go                     # pgxpool & schema-scoped goose migration runner
│   │   └── tenant_tx.go                    # Enforces tenant_id scoping on transactions
│   ├── events/
│   │   └── bus.go                          # Event struct, Bus interface, thread-safe InMemoryBus
│   └── httpx/
│       ├── response.go                     # Generic BindJSON[T](), JSON(), ProblemDetails Error()
│       ├── middleware.go                   # Gin Keycloak JWKS verify, RequireRole, TenantID, slog
│       └── auth_proxy.go                   # Headless /api/auth/login, /refresh, /logout via golang.org/x/oauth2
│
├── modules/                                # Isolated Go child modules (NO cross-module imports)
│   ├── organization/                       # Module 1: Org settings, punch policies, Keycloak user admin
│   │   ├── go.mod                          # "attendance_management/modules/organization"
│   │   ├── migrations/001_init.sql         # Schema: `org` (policies, employees)
│   │   ├── domain.go                       # AllowedPunchMethods, OrgPolicy structs
│   │   ├── repository.go                   # pgxpool queries scoped by tenant_id
│   │   ├── keycloak_admin.go               # Headless Keycloak Admin REST API client
│   │   ├── service.go                      # Policy enforcement & employee management
│   │   └── module.go                       # Gin handlers (/api/organization/*) & init()
│   │
│   ├── attendance/                         # Module 2: Punch-in/out logs, sessions, reports
│   │   ├── go.mod                          # "attendance_management/modules/attendance"
│   │   ├── migrations/001_init.sql         # Schema: `attendance` (logs, paired_sessions, tenant_policies)
│   │   ├── domain.go                       # PunchType (IN/OUT), PunchMethod, Session duration
│   │   ├── repository.go                   # Atomic punch-in/out queries via pgxpool
│   │   ├── service.go                      # Punch state machine, policy sync, report queries
│   │   └── module.go                       # Gin handlers (/api/attendance/*) & init()
│   │
│   └── payroll/                            # Module 3: Conditional salary slip generation
│       ├── go.mod                          # "attendance_management/modules/payroll"
│       ├── migrations/001_init.sql         # Schema: `payroll` (tenant_settings, salary_structures, salary_slips)
│       ├── domain.go                       # SalarySlip calculation & eligibility rules
│       ├── repository.go                   # pgxpool queries scoped by tenant_id
│       ├── service.go                      # Generates slips if org enabled; listens to events
│       └── module.go                       # Gin handlers (/api/payroll/*) & init()
│
│   └── user/                               # Module 4: User details and preferences management
│       ├── go.mod                          # "attendance_management/modules/user"
│       ├── migrations/001_init.sql         # Schema: `user` (users, preferences)
│       ├── domain.go                       # User roles and preference structures
│       ├── repository.go                   # pgxpool queries scoped by tenant_id
│       ├── service.go                      # User management and preferences handling
│       └── module.go                       # Gin handlers (/api/user/*) & init()
│
├── frontend-web/                           # Angular + Angular Material UI
│   └── src/app/
│       ├── core/                           # Native Material login, /api/auth/* token service, role guards
│       ├── shared/                         # Material tables, date-range pickers, status badges
│       └── features/
│           ├── employee-portal/            # Web Punch-In/Out, personal logs, salary slips
│           └── org-admin/                  # Punch method toggles, attendance reports, payroll
│
└── mobile-flutter/                         # Flutter + Material 3 UI (Dual-Mode App)
    └── lib/
        ├── core/                           # Native Material 3 login, /api/auth/* token storage, HTTP client
        ├── models/                         # Dart DTOs matching Go module contracts
        ├── services/                       # Attendance, OrgPolicy, Biometric hardware services
        └── features/
            ├── employee/                   # Mobile Punch-In/Out, history, salary slip viewer
            └── gate_kiosk/                 # Gate Terminal UI (FaceScan & FingerScan modes)
```