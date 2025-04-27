# Plan to make Gogrok look ok

UI/UX Improvement Plan (Detailed & Concrete)
Let's make Gogrok more visually appealing and user-friendly.

Overall Layout & Framework:

Adopt a Modern CSS Framework: Integrate Tailwind CSS or Bootstrap 5. This provides a consistent design system (spacing, typography, colors, components) and makes responsiveness easier. Tailwind is utility-first and highly customizable; Bootstrap is component-based. Let's assume Tailwind CSS for this plan (requires a build step for CSS).
Layout Structure:
Header: Simple header with the "Gogrok" title/logo.
Sidebar (Persistent): Main navigation. Items: Dashboard, Local Modules, External Modules, Security Overview (New!), Settings (Future). Use icons alongside text.
Main Content Area: This is where HTMX will swap content based on sidebar navigation or interactions within the content. Include Breadcrumbs at the top of this area to show context (e.g., Local Modules > github.com/org/module > pkg/subpkg).
Homepage / Dashboard (New / route):

Replace Initial Module List: Instead of just listing local modules, show a dashboard with summary cards/widgets:
"Total Local Modules": Count.
"Total External Dependencies": Count.
"Total Lines of Code (Local)": Sum across local modules.
"Avg. Complexity (Local)": Average across local modules.
"Modules with Security Issues": Count (Link to Security Overview).
"Outdated Dependencies": Count (Requires version comparison logic).
Maybe Charts: A simple bar chart showing LoC distribution per module, or a pie chart of dependency types. (Use a JS charting library like Chart.js, perhaps loaded via a simple HTMX extension or vanilla JS).
Module Lists (Local & External):

Use Tables: Display modules in a sortable table (<table> structure styled with Tailwind).
Columns: Module Path, Latest Version (Local)/Versions Used (External), LoC, Avg. Complexity, Direct Deps (Count), Reverse Deps (Count), Security Issues (Icon/Count).
Visual Cues: Use icons/colors: Red shield for vulnerabilities, yellow warning sign for outdated dependencies.
Filtering/Search: Add a search input above the table that filters the list dynamically (can be done with HTMX hx-trigger="keyup changed delay:500ms" pointing to the list endpoint with a search query param).
Module Detail View (/module/{path}):

Tabbed Interface: Use tabs (styled links controlling HTMX content swaps within a dedicated area) for different sections:
Overview: Show key stats (LoC, Complexity, Files, Packages), versions used, latest local tag/latest remote version.
Dependencies: List direct dependencies (table format: Module Path, Version Constraint). Add a Visualize button.
Reverse Dependencies: List local modules/packages that depend on this module (table format: Dependent Module/Package Path).
Packages: List packages within this module (table format similar to Module List: Name, LoC, Complexity, Files, Security Issues within package). Make package names clickable to load Package Detail View.
Security: Display detailed vulnerability information affecting this module or its dependencies (as determined by govulncheck/API).
(Future) Code Metrics: Show Maintainability Index, Coupling info, etc.
Dependency Visualization: Clicking "Visualize" could trigger an HTMX request that returns a div containing data for Mermaid.js or Cytoscape.js, which then renders the graph client-side. Show immediate dependencies initially, maybe allow expanding levels.
Package Detail View (/package/{path}?package=...):

Similar Structure: Use tabs or sections.
Overview: Stats (LoC, Files, Complexity, % Generated).
Files: List files (Table: Name, LoC, Complexity, Type). Make file names clickable.
Reverse Dependencies: List packages (local) that import this package (Table: Module Path, Package Name).
Imports: List packages this package imports (useful for seeing its direct external surface area).
File View (/file/{path}?package=...&file=...):

Syntax Highlighting: Load a JS syntax highlighting library (like Prism.js or Highlight.js). After HTMX loads the file content into a <pre><code> block, trigger the highlighting JS function. Ensure the Go language is supported/loaded.
Line Numbers: Add line numbers using CSS or the highlighting library's features.
Context: Clearly show Module, Package, and File name. Use breadcrumbs.
HTMX & Interactivity:

Ensure smooth transitions (e.g., add loading indicators using HTMX CSS classes).
Use hx-push-url="true" where appropriate to update the browser URL bar as the user navigates.
Leverage events for things like triggering JS (syntax highlighting, graph rendering) after content is swapped (htmx:afterSwap).
Styling Details:

Choose a clean color palette (e.g., based on Tailwind's defaults or a custom one).
Use consistent typography.
Employ icons (e.g., from Heroicons - integrates well with Tailwind) for navigation, buttons, and status indicators.
