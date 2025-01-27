document.addEventListener("DOMContentLoaded", () => {
  // Load navbar
  fetch("/static/components/navbar.html")
    .then((response) => response.text())
    .then((html) => {
      document.getElementById("navbar-placeholder").innerHTML = html;
      initializeNavbar();
    });

  // Theme management
  const theme = localStorage.getItem("theme") || "dark";
  document.documentElement.setAttribute("data-theme", theme);
});

function initializeNavbar() {
  const mobileMenuBtn = document.querySelector(".mobile-menu-btn");
  const navLinks = document.querySelector(".nav-links");
  const themeToggle = document.querySelector(".theme-toggle");

  // Toggle mobile menu
  mobileMenuBtn?.addEventListener("click", () => {
    navLinks.classList.toggle("mobile-visible");
  });

  // Toggle theme
  themeToggle?.addEventListener("click", () => {
    const currentTheme = document.documentElement.getAttribute("data-theme");
    const newTheme = currentTheme === "dark" ? "light" : "dark";
    document.documentElement.setAttribute("data-theme", newTheme);
    localStorage.setItem("theme", newTheme);
    themeToggle.textContent = newTheme === "dark" ? "🌙" : "🔆";
  });

  // Check authentication status and update navbar
  checkAuthStatus().then((isAuthenticated) => {
    updateAuthVisibility(isAuthenticated);
  });
}

function updateAuthVisibility(isAuthenticated) {
  const authHideElements = document.querySelectorAll(".auth-hide");
  const authShowElements = document.querySelectorAll(".auth-show");
  const authRequiredElements = document.querySelectorAll(".auth-required");

  if (isAuthenticated) {
    // User is signed in
    authHideElements.forEach((element) => {
      element.style.display = "none";
    });
    authShowElements.forEach((element) => {
      element.style.display = "block";
    });
    authRequiredElements.forEach((element) => {
      element.style.display = "block";
    });
  } else {
    // User is not signed in
    authHideElements.forEach((element) => {
      element.style.display = "block";
    });
    authShowElements.forEach((element) => {
      element.style.display = "none";
    });
    authRequiredElements.forEach((element) => {
      element.style.display = "none";
    });
  }
}

async function checkAuthStatus() {
  try {
    const response = await fetch("/api/auth/status");
    if (!response.ok) {
      throw new Error("Failed to fetch authentication status");
    }
    const data = await response.json();
    return data.isAuthenticated;
  } catch (error) {
    console.error("Error checking auth status:", error);
    return false;
  }
}
