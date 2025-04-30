document.addEventListener("DOMContentLoaded", () => {
    const profileButton = document.getElementById("profile-button");
    const userDropdown = document.getElementById("user-dropdown");
  
    profileButton.addEventListener("click", () => {
      userDropdown.classList.toggle("hidden");
    });
  
    // Optional: Close the dropdown when clicking outside
    document.addEventListener("click", (event) => {
      if (!profileButton.contains(event.target) && !userDropdown.contains(event.target)) {
        userDropdown.classList.add("hidden");
      }
    });
  });