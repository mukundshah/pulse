# Frontend Implementation Plan - Checkly Design

## Overview

Rebuild all frontend pages from scratch, implementing a complete monitoring dashboard matching Checkly's design patterns. Use full-page forms for check creation/editing, implement comprehensive metrics visualization, and create a sidebar-based navigation structure.

## Architecture Decisions

### UI Patterns (Based on Checkly Design)

- **Full-Page Forms**: Use for check creation/editing (not dialogs) - matches Checkly's multi-section form layout
- **Sidebar Navigation**: Keep existing sidebar structure (Dashboard, Projects, Settings, etc.)
- **Status Overview Cards**: Large status cards at top (PASSING, DEGRADED, FAILING) with counts
- **Tables with Sparklines**: Check lists with inline sparkline graphs for recent activity
- **Right Sidebar**: Use for run results/details in check detail pages
- **Metrics Cards**: Display key metrics (P50, P95, P99, Availability, etc.) with percentage changes
- **Time Range Filters**: Prominent time range buttons (Today, 1hr, 3hr, 24hr, 7d, Custom)
- **Status Filters**: Filter buttons for Passed, Failed, Degraded, Has retries, Location
- **Graphs/Charts**: Multiple chart types (line graphs, bar charts) for performance visualization
- **Dialogs**: Use sparingly for quick actions (invite member, change role)

### Component Structure

- **Shared Components**: Create reusable components in `web/components/`
  - `StatusOverviewCards.vue` - Large status cards (PASSING, DEGRADED, FAILING) with counts
  - `TimeRangeFilters.vue` - Time range buttons (Today, 1hr, 3hr, 24hr, 7d, Custom)
  - `StatusFilters.vue` - Status filter buttons (Passed, Failed, Degraded, Has retries, Location)
  - `SparklineChart.vue` - Inline sparkline graphs for check activity
  - `MetricsCard.vue` - Metric display with value, percentage change, and color coding
  - `PerformanceGraph.vue` - Line/bar charts for performance metrics
  - `CheckFormHttp.vue` - HTTP check form component
  - `CheckFormTcp.vue` - TCP check form component
  - `CheckFormDns.vue` - DNS check form component
  - `CheckFormBrowser.vue` - Browser check form component
  - `CheckFormHeartbeat.vue` - Heartbeat check form component
  - `LocationSelector.vue` - Location selection grid (public/private locations)
  - `AssertionsTable.vue` - Assertions configuration table
  - `RunResultsSidebar.vue` - Right sidebar for check run results
  - `AlertCard.vue` - Alert display component
  - `StatusBadge.vue` - Status indicator component
  - `MemberInviteDialog.vue` - Invite member dialog (small dialog)
  - `MemberRoleDialog.vue` - Change member role dialog (small dialog)

### API Integration

- Use `useAPI` composable from `web/modules/api/runtime/composable.ts` or `$api` from `useNuxtApp()` directly in components/pages where needed
- No need to create wrapper composables - consume the API directly

## Page-by-Page Implementation

### 1. Layout Updates (`web/layouts/app.vue`)

**Keep Existing Sidebar Structure:**

- Fix icon typo: `lucide:pluse` → `lucide:plus` (line 49)
- Keep the current sidebar structure (Dashboard, Projects with collapsible items, Settings, Integrations)
- Make SidebarGroupAction clickable to show inline project creation input (for Projects section at line 48-50)
  - When clicked, display an input field directly within the sidebar below the "Projects" label
  - Input should have placeholder "Project name" and submit on Enter
  - Show cancel/close button to dismiss the input
  - On successful creation, add the new project to the sidebar list and clear the input
- Dynamically load projects from API and populate sidebar
- Make sidebar project items clickable and show active project state
- Update sidebar links to use dynamic project IDs
- Keep existing sidebar navigation items (Dashboard, Checks, Alerts, Status, Settings under each project)

### 2. Dashboard/Home (`web/pages/dashboard.vue` or `web/pages/index.vue`)

**Layout (Matching Checkly Home):**

- **Status Overview Cards**: Three large cards at top
  - PASSING (green) with count
  - DEGRADED (orange/yellow) with count
  - FAILING (red) with count
- **Search and Filter Bar**:
  - Search input: "Search by name, request url..."
  - Filter dropdowns: Status, Check type, Tags
- **Checks Table**:
  - Columns: NAME, TYPE, LAST RESULTS (sparkline), 24H, 7D, AVG, P95, ΔT
  - Each row shows: Check name with status icon, type badge, sparkline graph, uptime percentages, response times
  - Last check time below name
  - Actions menu (ellipsis) on right

**Components:**

- StatusOverviewCards component
- SearchBar with filters
- ChecksTable with sparklines
- SparklineChart component

**API Calls:**

- GET /projects/:projectId/checks - List all checks
- Aggregate status counts
- Calculate uptime percentages and response times

### 3. Project Detail (`web/pages/projects/[id]/index.vue`)

**Layout:**

- Same as Dashboard/Home layout (status cards, search, checks table)
- Or simplified view showing project-specific checks
- "New Check" button in header (opens full-page form)

**Components:**

- Same as Dashboard components
- Check form components (type-specific, full-page, not dialog)

**API Calls:**

- GET /projects/:id - Project details
- GET /projects/:id/checks - List checks
- POST /projects/:id/checks - Create check
- PUT /projects/:id - Update project (if needed)

### 4. Check Detail (`web/pages/projects/[id]/checks/[checkId].vue`)

**Layout (Matching Checkly Check Detail):**

- **Header**:
  - Check name with status icon (green circle for passing)
  - Location flags (Australia, Indonesia, New Zealand, Singapore)
  - Large status indicator: "Check is passing" with green checkmark
  - Last updated time and "Edit" button
- **Time Range Filters**: Custom, Today, 1hr, 3hr, 24hr, 7d buttons
- **Status Filters**: Passed (green), Failed (red), Degraded (orange), Has retries, Location dropdown
- **Key Metrics Cards**:
  - Availability: 100% with percentage change
  - Retries: 0% with percentage change
  - P50: 2.74s with percentage change
  - P95: 3.9s with percentage change
  - Failure Alerts: 0
  - Span Errors: 0
- **Historical Graph**: Bar chart showing check results over time
- **Performance Section** (if applicable):
  - Check duration graph (P50, P95, P99 lines)
  - Loading metrics (TTFB, FCP, LCP, Loaded, DCL) with graph
  - Errors section (Console, Network, Script, Document) with bar chart
  - Interactivity (TBT) with line graph
- **Error Groups Table**: MESSAGE, FIRST SEEN, LAST SEEN, EVENTS, LOCATIONS
- **Alerts Table**: STATUS, LOCATION, NOTIFICATIONS, TIMESTAMP
- **Locations Table**: LOCATION, P50, P95, P99
- **Right Sidebar**: Run results list with location, duration, timestamp

**Components:**

- TimeRangeFilters component
- StatusFilters component
- MetricsCard components
- PerformanceGraph components (line and bar charts)
- RunResultsSidebar component
- ErrorGroupsTable component
- AlertsTable component
- LocationsTable component

**API Calls:**

- GET /checks/:checkId - Check details
- GET /checks/:checkId/runs?limit=100 - Check runs (for sidebar and graphs)
- GET /checks/:checkId/alerts?limit=100 - Check alerts
- PUT /checks/:checkId - Update check
- DELETE /checks/:checkId - Delete check

### 5. Check Form (`web/pages/projects/[id]/checks/new/[type].vue` and `web/pages/projects/[id]/checks/[checkId]/edit.vue`)

**Route Structure:**

- `new/[type].vue` where `type` can be: `http`, `tcp`, `dns`, `browser`, `heartbeat`
- The page component will render the appropriate form component based on the type parameter
- For edit, use `[checkId]/edit.vue` which determines the type from the check data

**Layout (Matching Checkly New URL Monitor):**

- **Header**: "New [type] monitor" or "Edit [type] monitor" with back arrow and "Create monitor" button (Ctrl+S)
- **Common Sections** (shared across all types):
  - **Section 1: Monitor Name & Tags**
    - Large title input: "[Type] Monitor #2"
    - Tag input: "Type a tag, hit enter" with location flags
    - Checkboxes: "Activated" (checked), "Muted" (unchecked) with info icons
  - **Section 4: Response Time Limits**
    - "DEGRADED AFTER": Input (3000) with Milliseconds/Seconds toggle
    - "FAILED AFTER": Input (5000) with Milliseconds/Seconds toggle
    - Note about 30-second timeout cap
  - **Section 5: Scheduling Strategy**
    - Radio buttons: "Parallel runs" (ENTERPRISE, disabled), "Round-robin" (selected)
    - Description text
  - **Section 6: Locations**
    - Private locations section (TEAM tag, upgrade button if not available)
    - Public locations section:
      - Description about selecting 2-3 locations
      - "2 locations selected" info box
      - "Clear selection" link
      - Three-column grid: Americas, Europe/Middle East/Africa, Asia Pacific
      - Location cards with flag, name, region code, TEAM tags
      - Selected locations highlighted in blue

**Type-Specific Sections:**

- **HTTP Check Form** (`CheckFormHttp.vue`):
  - **Section 2: Monitor a URL**
    - Description text with link to environment variables
    - IPv4/IPv6 dropdown
    - URL input with method selector (GET, POST, etc.) on left
    - URL field with pre-filled example
    - Checkboxes: "Skip SSL", "Follow redirects" (checked), "This request should fail"
    - Location dropdown (flag icon)
    - "Run request" button (Ctrl+→)
  - **Section 3: Assertions**
    - Description with link to documentation
    - Table: SOURCE (Status code dropdown), PROPERTY (empty), COMPARISON (Equals dropdown), TARGET (input)
    - "Add" button to add more assertions

- **TCP Check Form** (`CheckFormTcp.vue`):
  - Host and port inputs
  - Connection timeout settings

- **DNS Check Form** (`CheckFormDns.vue`):
  - Domain name input
  - Record type selector (A, AAAA, MX, etc.)
  - Expected value validation

- **Browser Check Form** (`CheckFormBrowser.vue`):
  - URL input
  - Playwright script editor
  - Screenshot options

- **Heartbeat Check Form** (`CheckFormHeartbeat.vue`):
  - Simple configuration (minimal fields)

**Components:**

- CheckFormHeader component (shared)
- MonitorNameSection component (shared)
- ResponseTimeLimitsSection component (shared)
- SchedulingStrategySection component (shared)
- LocationSelector component (shared)
- CheckFormHttp component (type-specific)
- CheckFormTcp component (type-specific)
- CheckFormDns component (type-specific)
- CheckFormBrowser component (type-specific)
- CheckFormHeartbeat component (type-specific)
- AssertionsTable component (for HTTP)

**Validation:**

- Name required (all types)
- Type-specific validation:
  - HTTP: URL required and valid, assertions validation
  - TCP: Host and port required
  - DNS: Domain required, record type required
  - Browser: URL required, script validation
  - Heartbeat: Minimal validation
- At least 2 locations for public locations (all types)
- Response time limits validation (all types)

### 6. Alerts Page (`web/pages/alerts.vue`)

**Layout:**

- Header with title and filter buttons (All, Active, Resolved)
- Stats cards: Total Alerts, Active, Resolved
- Alerts table: Status, Severity, Check Name, Message, Triggered At, Resolved At, Actions

**Components:**

- AlertCard component
- AlertFilters component
- Status badges

**API Calls:**

- Aggregate alerts from all checks (may need to fetch per check or implement backend aggregation)
- Filter by status (active/resolved)

### 7. Status Page (`web/pages/status.vue`)

**Layout:**

- Large status banner at top (Operational/Degraded/Major Outage)
- Services grid: List of all checks/services with status
- Recent Updates timeline: Past incidents and maintenance

**Components:**

- StatusBanner component
- ServiceStatusCard components
- UpdateTimeline component

**API Calls:**

- Aggregate status from all checks
- Calculate overall status based on check statuses

### 8. Project Settings (`web/pages/projects/[id]/settings.vue`)

**Layout:**

- Sub-tabs: Status Page, Notifications, Members
- Status Page tab: URL display, embed code, view/configure buttons
- Notifications tab: Webhooks table with name, URL, events, status, success rate, actions
- Members tab: Members table with avatar, name, email, role, joined date, actions

**Components:**

- SettingsTabs component
- WebhooksTable component
- WebhookForm dialog
- MembersTable component
- MemberInviteDialog component
- MemberRoleDialog component

**API Calls:**

- GET /projects/:id/members - List members
- POST /projects/:id/invites - Create invite
- GET /projects/:id/invites - List invites
- PUT /projects/:id/members/:userId - Update role
- DELETE /projects/:id/members/:userId - Remove member

## Implementation Details

### Data Fetching Strategy

1. **Initial Load**: Fetch data on page mount
2. **Caching**: Use Nuxt's built-in caching for GET requests
3. **Refetch**: Manual refresh buttons where needed
4. **Optimistic Updates**: For create/update/delete operations
5. **Error Handling**: Show error toasts, handle 404s gracefully

### Status Indicators

- **Passing/Operational**: Green dot, "Passing" label
- **Degraded/Warning**: Orange/yellow dot, "Degraded" label
- **Failing/Critical**: Red dot, "Failing" label
- **Unknown**: Gray dot, "Unknown" label

### Time Formatting

- Use relative time for recent items ("2 minutes ago", "about 16 hours ago")
- Use absolute time for older items
- Always use NuxtTime component

### Chart Libraries

- Use a charting library (e.g., Chart.js, Recharts, or similar) for:
  - Sparkline charts in tables
  - Line graphs for performance metrics
  - Bar charts for error counts
  - Historical check results visualization

### Empty States

- Projects: "No projects yet. Create your first project to get started."
- Checks: "No checks configured. Create a check to start monitoring."
- Alerts: "No alerts. All systems operational."
- Check Runs: "No runs yet. Check will run on next scheduled interval."

### Loading States

- Skeleton loaders for tables and cards
- Spinner for buttons during actions
- Progress indicators for long operations

### Error Handling

- Toast notifications for errors
- Inline validation for forms
- 404 pages for missing resources
- Retry mechanisms for failed API calls

## File Structure

```
web/
├── components/
│   ├── projects/
│   │   ├── ProjectCard.vue
│   │   └── ProjectHeader.vue
│   ├── checks/
│   │   ├── CheckFormHttp.vue
│   │   ├── CheckFormTcp.vue
│   │   ├── CheckFormDns.vue
│   │   ├── CheckFormBrowser.vue
│   │   ├── CheckFormHeartbeat.vue
│   │   ├── ChecksTable.vue (with sparklines)
│   │   ├── CheckRunsTable.vue
│   │   ├── RunResultsSidebar.vue
│   │   ├── AssertionsTable.vue
│   │   ├── LocationSelector.vue
│   │   └── SchedulingStrategySection.vue
│   ├── metrics/
│   │   ├── StatusOverviewCards.vue
│   │   ├── MetricsCard.vue
│   │   ├── PerformanceGraph.vue
│   │   └── SparklineChart.vue
│   ├── filters/
│   │   ├── TimeRangeFilters.vue
│   │   ├── StatusFilters.vue
│   │   └── SearchBar.vue
│   ├── alerts/
│   │   ├── AlertCard.vue
│   │   └── AlertsTable.vue
│   ├── members/
│   │   ├── MemberInviteDialog.vue
│   │   ├── MemberRoleDialog.vue
│   │   └── MembersTable.vue
│   ├── status/
│   │   ├── StatusBanner.vue
│   │   ├── ServiceStatusCard.vue
│   │   └── UpdateTimeline.vue
│   └── shared/
│       ├── StatusBadge.vue
│       ├── EmptyState.vue
│       └── LocationFlag.vue
└── pages/
    ├── dashboard.vue (rebuild - Checkly home layout)
    ├── alerts.vue (rebuild)
    ├── status.vue (rebuild)
    └── projects/
        └── [id]/
            ├── index.vue (rebuild - same as dashboard but project-specific)
            ├── settings.vue (rebuild)
            └── checks/
                ├── new/
                │   └── [type].vue (renders appropriate form component based on type)
                ├── [checkId].vue (rebuild - detailed view with graphs)
                └── [checkId]/edit.vue (renders appropriate form component based on check type)
```

## Implementation Order

1. **Setup & Infrastructure**

   - Create shared components (StatusBadge, EmptyState, LocationFlag)
   - Fix sidebar icon and implement project creation trigger
   - Set up chart library

2. **Core Components**

   - StatusOverviewCards component
   - SparklineChart component
   - MetricsCard component
   - TimeRangeFilters component
   - StatusFilters component

3. **Core Pages**

   - Dashboard/Home page (with status cards and checks table)
   - Project detail page
   - Check detail page (with graphs and metrics)

4. **Forms**

   - Check form components (type-specific: HTTP, TCP, DNS, Browser, Heartbeat)
   - Location selector component
   - Assertions table component (for HTTP checks)

5. **Supporting Pages**

   - Alerts page
   - Status page

6. **Settings & Management**

   - Project settings page
   - Member management
   - Invite management

7. **Polish & Refinement**

   - Loading states
   - Error handling
   - Empty states
   - Responsive design
   - Accessibility improvements
   - Chart animations and interactions
