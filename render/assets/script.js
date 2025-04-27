document.body.addEventListener('htmx:afterRequest', function (evt) {
    const errorTarget = document.getElementById("htmx-alert")
    if (evt.detail.successful) {
        // Successful request, clear out alert
        errorTarget.setAttribute("hidden", "true")
        errorTarget.innerText = "";

        // Update breadcrumbs based on the requested path
        updateBreadcrumbs(evt.detail.pathInfo.requestPath);
    } else if (evt.detail.failed && evt.detail.xhr) {
        // Server error with response contents, equivalent to htmx:responseError
        console.warn("Server error", evt.detail)
        const xhr = evt.detail.xhr;
        errorTarget.innerText = `Unexpected server error: ${xhr.status} - ${xhr.statusText}`;
        errorTarget.removeAttribute("hidden");
    } else {
        // Unspecified failure, usually caused by network error
        console.error("Unexpected htmx error", evt.detail)
        errorTarget.innerText = "Unexpected error, check your connection and try to refresh the page.";
        errorTarget.removeAttribute("hidden");
    }
});

// Function to update breadcrumbs based on the current path
function updateBreadcrumbs(path) {
    const breadcrumbCurrent = document.getElementById("breadcrumb-current");
    if (!breadcrumbCurrent) return;

    // Default to Dashboard
    let currentPage = "Dashboard";

    // Extract the page name from the path
    if (path === "/dashboard") {
        currentPage = "Dashboard";
    } else if (path === "/local") {
        currentPage = "Local Modules";
    } else if (path === "/external") {
        currentPage = "External Modules";
    } else if (path === "/about") {
        currentPage = "About";
    } else if (path.startsWith("/module/")) {
        currentPage = "Module Details";
    } else if (path.startsWith("/package/")) {
        currentPage = "Package Details";
    } else if (path.startsWith("/file/")) {
        currentPage = "File Details";
    }

    breadcrumbCurrent.textContent = currentPage;
}

// Initialize breadcrumbs on page load
document.addEventListener('DOMContentLoaded', function() {
    // Set initial breadcrumb to Dashboard
    updateBreadcrumbs("/dashboard");

    // Highlight the active sidebar item
    const sidebarLinks = document.querySelectorAll('aside nav a');
    sidebarLinks.forEach(link => {
        link.addEventListener('click', function() {
            // Remove active class from all links
            sidebarLinks.forEach(l => {
                l.classList.remove('text-gray-900', 'bg-gray-100');
                l.classList.add('text-gray-600', 'hover:text-gray-900', 'hover:bg-gray-50');
            });

            // Add active class to clicked link
            this.classList.remove('text-gray-600', 'hover:text-gray-900', 'hover:bg-gray-50');
            this.classList.add('text-gray-900', 'bg-gray-100');
        });
    });
});
