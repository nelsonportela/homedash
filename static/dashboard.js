// Health check async logic with progress polling
window.addEventListener('DOMContentLoaded', () => {
  const healthLoading = document.getElementById('healthLoading');
  const healthProgressText = document.getElementById('healthProgressText');
  const healthProgressBar = document.getElementById('healthProgressBar');
  const servicesDownGroup = document.getElementById('servicesDownGroup');
  const servicesDownApps = document.getElementById('servicesDownApps');
  if (!healthLoading) return;

  let pollInterval;
  let total = 0;
  let checked = 0;
  let finished = false;

  function updateProgressBar(checked, total) {
    if (total === 0) return;
    const percent = Math.round((checked / total) * 100);
    healthProgressBar.style.width = percent + '%';
    healthProgressText.textContent = `${checked} / ${total} services checked...`;
  }

  function pollProgress() {
    fetch('/api/health/progress')
      .then(res => res.json())
      .then(data => {
        if (data && typeof data.checked === 'number' && typeof data.total === 'number') {
          checked = data.checked;
          total = data.total;
          updateProgressBar(checked, total);
        }
      });
  }

  // Start polling progress
  pollInterval = setInterval(pollProgress, 250);
  pollProgress();

  fetch('/api/health')
    .then(res => res.json())
    .then(data => {
      finished = true;
      clearInterval(pollInterval);
      healthLoading.style.display = 'none';
      if (data && data.down && Object.keys(data.down).length > 0) {
        // Find all app shortcuts
        const allApps = document.querySelectorAll('.app-item');
        let found = 0;
        allApps.forEach(app => {
          const appName = app.getAttribute('data-app-name');
          if (data.down[appName]) {
            // Move to Services Down group
            servicesDownApps.appendChild(app);
            found++;
          }
        });
        if (found > 0) {
          servicesDownGroup.classList.remove('hidden');
        }
      } else {
        servicesDownGroup.classList.add('hidden');
      }
    })
    .catch(() => {
      finished = true;
      clearInterval(pollInterval);
      healthLoading.style.display = 'none';
      // Optionally show error
    });
});
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