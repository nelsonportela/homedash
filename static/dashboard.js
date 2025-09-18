// Theme management
const themeToggle = document.getElementById('themeToggle');
const sunIcon = document.getElementById('sunIcon');
const moonIcon = document.getElementById('moonIcon');
const body = document.body;
const html = document.documentElement;

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme');
  const defaultTheme = body.getAttribute('data-default-theme') || 'dark';
  const theme = savedTheme || defaultTheme;
  
  if (theme === 'dark') {
    html.classList.add('dark');
    body.classList.add('dark');
    sunIcon.classList.remove('hidden');
    moonIcon.classList.add('hidden');
  } else {
    html.classList.remove('dark');
    body.classList.remove('dark');
    sunIcon.classList.add('hidden');
    moonIcon.classList.remove('hidden');
  }
}

// Toggle theme
themeToggle.addEventListener('click', () => {
  const isDark = html.classList.contains('dark');
  
  if (isDark) {
    html.classList.remove('dark');
    body.classList.remove('dark');
    sunIcon.classList.add('hidden');
    moonIcon.classList.remove('hidden');
    localStorage.setItem('theme', 'light');
  } else {
    html.classList.add('dark');
    body.classList.add('dark');
    sunIcon.classList.remove('hidden');
    moonIcon.classList.add('hidden');
    localStorage.setItem('theme', 'dark');
  }
});

// Search functionality
const searchInput = document.getElementById('searchInput');
const appGroups = document.getElementById('appGroups');
const noResults = document.getElementById('noResults');

searchInput.addEventListener('input', (e) => {
  const searchTerm = e.target.value.toLowerCase().trim();
  const groups = document.querySelectorAll('.app-group');
  let hasVisibleApps = false;
  
  groups.forEach(group => {
    const apps = group.querySelectorAll('.app-item');
    let hasVisibleAppsInGroup = false;
    
    apps.forEach(app => {
      const appName = app.getAttribute('data-app-name').toLowerCase();
      const appTitle = app.querySelector('.text-lg').textContent.toLowerCase();
      const isVisible = appName.includes(searchTerm) || appTitle.includes(searchTerm);
      
      app.style.display = isVisible ? 'flex' : 'none';
      if (isVisible) {
        hasVisibleAppsInGroup = true;
        hasVisibleApps = true;
      }
    });
    
    group.style.display = hasVisibleAppsInGroup ? 'block' : 'none';
  });
  
  appGroups.style.display = hasVisibleApps ? 'block' : 'none';
  noResults.style.display = hasVisibleApps ? 'none' : 'block';
});

// Refresh functionality
const refreshBtn = document.getElementById('refreshBtn');
const refreshIcon = document.getElementById('refreshIcon');

refreshBtn.addEventListener('click', async () => {
  // Add spinning animation
  refreshIcon.classList.add('animate-spin');
  refreshBtn.disabled = true;
  
  try {
    const response = await fetch('/refresh', { method: 'POST' });
    const result = await response.json();
    
    if (response.ok && result.status === 'success') {
      // Show success feedback briefly before reload
      const originalText = refreshBtn.querySelector('span').textContent;
      refreshBtn.querySelector('span').textContent = 'Success!';
      
      // Reload the page after a short delay
      setTimeout(() => {
        window.location.reload();
      }, 800);
    } else {
      throw new Error(result.message || 'Refresh failed');
    }
  } catch (error) {
    console.error('Refresh error:', error);
    // Remove spinning animation on error
    refreshIcon.classList.remove('animate-spin');
    refreshBtn.disabled = false;
    
    // Show error feedback
    const originalText = refreshBtn.querySelector('span').textContent;
    refreshBtn.querySelector('span').textContent = 'Error';
    refreshBtn.title = error.message || 'Refresh failed';
    setTimeout(() => {
      refreshBtn.querySelector('span').textContent = originalText;
      refreshBtn.title = 'Refresh dashboard';
    }, 3000);
  }
});

// Initialize theme on page load
initTheme();