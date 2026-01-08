package main

import (
	 "crypto/md5"
		 "fmt"
		 "html/template"
		 "io"
		 "io/ioutil"
		 "log"
		 "net/http"
		 "os"
		 "path/filepath"
		 "regexp"
		 "sort"
		 "strings"
		 "gopkg.in/yaml.v3"
	 )

const Version = "v0.5.0"

// (health check functionality removed)

type App struct {
	Name    string
	Title   string
	URL     string
	Icon    string
	IconURL string // Local icon URL path
	Group   string
}

type AppConfig struct {
	IconURL   string `yaml:"icon_url"`
	Title     string `yaml:"title"`
	URL       string `yaml:"url"`  // Manual URL for standalone apps not in Caddyfile
	Show      *bool  `yaml:"show"` // Whether to show the app (defaults to true if not specified)
}

type AppGroup struct {
	Name     string
	Apps     []App
	Color    string // Group accent color
	GridCols string // Grid column class for this group
}

type DashboardData struct {
	Title      string
	Groups     []AppGroup
	DefaultTheme string // Default theme preference
}

type Config struct {
	Title         string                   `yaml:"title"`
	DefaultTheme  string                   `yaml:"default_theme"` // "dark" or "light"
	Groups        map[string]GroupConfig   `yaml:"groups"`        // Groups with their apps
	GroupOrder    []string                 // Preserves the order from YAML
}
// (health check helper removed)

type GroupConfig struct {
	Color    string                   `yaml:"color"`     // Hex color for group accent
	GridCols int                      `yaml:"grid_cols"` // Custom grid columns (2-8)
	Apps     map[string]AppConfig     `yaml:"apps"`      // Apps in this group
}

// Custom unmarshaler to preserve group order
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	// First unmarshal normally
	type ConfigAlias Config
	var config ConfigAlias
	if err := value.Decode(&config); err != nil {
		return err
	}
	
	// Extract group order from YAML manually
	var groupOrder []string
	
	// Handle both DocumentNode and MappingNode cases
	var targetNode *yaml.Node
	if value.Kind == yaml.DocumentNode && len(value.Content) > 0 {
		targetNode = value.Content[0]
	} else if value.Kind == yaml.MappingNode {
		targetNode = value
	}
	
	if targetNode != nil && targetNode.Kind == yaml.MappingNode {
		for i := 0; i < len(targetNode.Content); i += 2 {
			if targetNode.Content[i].Value == "groups" {
				groupsNode := targetNode.Content[i+1]
				if groupsNode.Kind == yaml.MappingNode {
					for j := 0; j < len(groupsNode.Content); j += 2 {
						groupOrder = append(groupOrder, groupsNode.Content[j].Value)
					}
				}
				break
			}
		}
	}
	
	*c = Config(config)
	c.GroupOrder = groupOrder
	return nil
}

func loadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return &Config{
			Groups:     make(map[string]GroupConfig),
			GroupOrder: []string{},
		}, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	if config.Groups == nil {
		config.Groups = make(map[string]GroupConfig)
	}
	if config.GroupOrder == nil {
		config.GroupOrder = []string{}
	}

	return &config, nil
}

func downloadIcon(url, name string) (string, error) {
	// Create cache directory
	cacheDir := "/tmp/icons"
	os.MkdirAll(cacheDir, 0755)
	
	// Generate filename from URL hash
	hash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	ext := filepath.Ext(url)
	if ext == "" {
		ext = ".png"
	}
	filename := name + "_" + hash[:8] + ext
	localPath := filepath.Join(cacheDir, filename)
	
	// Check if already downloaded
	if _, err := os.Stat(localPath); err == nil {
		return "/icons/" + filename, nil
	}
	
	// Download icon
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download icon: %d", resp.StatusCode)
	}
	
	// Save to file
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	
	return "/icons/" + filename, nil
}

func parseCaddyfile(path string, config *Config) ([]App, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Regex to find app blocks: @name host name.url
	appRe := regexp.MustCompile(`@([a-zA-Z0-9_-]+) host ([a-zA-Z0-9_.-]+)`)
	apps := []App{}

	for _, match := range appRe.FindAllStringSubmatch(string(content), -1) {
		name := match[1]
		url := "https://" + match[2]
		
		var icon string
		var iconURL string
		var displayName string
		var groupName string
		
		// Find which group this app belongs to (look through all groups and their apps)
		var appConfig *AppConfig
		for gName, groupConfig := range config.Groups {
			if ac, exists := groupConfig.Apps[name]; exists {
				appConfig = &ac
				groupName = gName
				break
			}
		}
		
		if appConfig != nil {
			// Check if app should be shown (defaults to true if not specified)
			if appConfig.Show != nil && !*appConfig.Show {
				continue // Skip this app if show is explicitly set to false
			}
			
			// Use custom title if provided, otherwise use default name
			if appConfig.Title != "" {
				displayName = appConfig.Title
			} else {
				displayName = strings.Title(name)
			}
			
			// Handle custom icon URL
			if appConfig.IconURL != "" {
				// Download and get local path
				if localPath, err := downloadIcon(appConfig.IconURL, name); err == nil {
					iconURL = localPath
				} else {
					log.Printf("Failed to download icon for %s: %v", name, err)
					icon = "🔗"
				}
			} else {
				// Use default icon
				icon = "🔗"
			}
		} else {
			// App not found in config, use defaults
			displayName = strings.Title(name)
			icon = "🔗"
			groupName = "Other"
		}
		
		apps = append(apps, App{
			Name:    name,        // The identifier from Caddyfile
			Title:   displayName, // The display name (custom title or default)
			URL:     url,
			Icon:    icon,
			IconURL: iconURL,
			Group:   groupName,
		})
	}
	return apps, nil
}

func getManualApps(config *Config) []App {
	var apps []App
	
	// Iterate through all groups and their apps
	for groupName, groupConfig := range config.Groups {
		for appName, appConfig := range groupConfig.Apps {
			// Only process entries that have a URL (manual entries)
			if appConfig.URL == "" {
				continue
			}
			
			// Check if app should be shown (defaults to true if not specified)
			if appConfig.Show != nil && !*appConfig.Show {
				continue
			}
			
			var icon string
			var iconURL string
			var displayName string
			
			// Use custom title if provided, otherwise use the config key name
			if appConfig.Title != "" {
				displayName = appConfig.Title
			} else {
				displayName = strings.Title(appName)
			}
			
			// Handle custom icon URL
			if appConfig.IconURL != "" {
				// Download and get local path
				if localPath, err := downloadIcon(appConfig.IconURL, appName); err == nil {
					iconURL = localPath
				} else {
					log.Printf("Failed to download icon for %s: %v", appName, err)
					icon = "🔗"
				}
			} else {
				// Use default icon
				icon = "🔗"
			}
			
			apps = append(apps, App{
				Name:    appName,      // The identifier from config
				Title:   displayName,  // The display name (custom title or default)
				URL:     appConfig.URL,
				Icon:    icon,
				IconURL: iconURL,
				Group:   groupName,    // Group is now taken from the parent group
			})
		}
	}
	
	return apps
}

func organizeAppsByGroup(apps []App, config *Config) []AppGroup {
	groupMap := make(map[string][]App)
	
	// Group apps by their group name
	for _, app := range apps {
		groupMap[app.Group] = append(groupMap[app.Group], app)
	}
	
	// Helper function to get grid class for group width (out of 12 columns)
	getGridClass := func(cols int) string {
		switch cols {
		case 1:
			return "1"
		case 2:
			return "2" 
		case 3:
			return "3"
		case 4:
			return "4"
		case 5:
			return "5"
		case 6:
			return "6"
		case 7:
			return "7"
		case 8:
			return "8"
		case 9:
			return "9"
		case 10:
			return "10"
		case 11:
			return "11"
		case 12:
			return "12"
		default:
			return "12" // Full width by default
		}
	}
	
	// Helper function to get display name for sorting
	getDisplayName := func(app App) string {
		if app.Title != "" {
			return app.Title
		}
		return app.Name
	}
	
	// Sort apps alphabetically within each group
	for groupName, apps := range groupMap {
		sort.Slice(apps, func(i, j int) bool {
			return strings.ToLower(getDisplayName(apps[i])) < strings.ToLower(getDisplayName(apps[j]))
		})
		groupMap[groupName] = apps // Update the map with sorted slice
	}
	
	// Convert map to slice using preserved order from config
	var groups []AppGroup
	
	// Add groups in the order they appear in the config file
	for _, groupName := range config.GroupOrder {
		if apps, exists := groupMap[groupName]; exists {
			groupConfig := config.Groups[groupName]
			
			group := AppGroup{
				Name:     groupName,
				Apps:     apps,
				Color:    groupConfig.Color,
				GridCols: getGridClass(groupConfig.GridCols),
			}
			
			groups = append(groups, group)
			delete(groupMap, groupName)
		}
	}
	
	// Add any remaining groups not in the preferred order
	for groupName, apps := range groupMap {
		groupConfig := config.Groups[groupName]
		
		group := AppGroup{
			Name:     groupName,
			Apps:     apps,
			Color:    groupConfig.Color,
			GridCols: getGridClass(groupConfig.GridCols),
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
       config, err := loadConfig("/etc/config/config.yaml")
       if err != nil {
	       log.Printf("Failed to load config: %v", err)
	       config = &Config{
		       Groups:     make(map[string]GroupConfig),
		       GroupOrder: []string{},
	       }
       }

       apps, err := parseCaddyfile("/etc/caddy/Caddyfile", config)
       if err != nil {
	       http.Error(w, "Failed to parse Caddyfile", 500)
	       return
       }

       // Get manual apps from config and merge them
       manualApps := getManualApps(config)
       allApps := append(apps, manualApps...)

       // Organize all apps by group (no health check)
       groups := organizeAppsByGroup(allApps, config)

       // Use configured title or default
       title := config.Title
       if title == "" {
	       title = "HomeDash"
       }

       // Use configured theme or default
       defaultTheme := config.DefaultTheme
       if defaultTheme == "" {
	       defaultTheme = "dark"
       }

       data := DashboardData{
	       Title:        title,
	       Groups:       groups,
	       DefaultTheme: defaultTheme,
       }

       // Load template from file
       tmpl, err := template.ParseFiles("templates/dashboard.html")
       if err != nil {
	       http.Error(w, "Template file not found: "+err.Error(), 500)
	       return
       }

       tmpl.Execute(w, data)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	
	// Try to reload the config to verify it's valid
	_, err := loadConfig("/etc/config/config.yaml")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"status": "error", "message": "Config reload failed: %v"}`, err)))
		return
	}
	
	// Try to parse the Caddyfile to verify it's valid
	config := &Config{
		Groups:     make(map[string]GroupConfig),
		GroupOrder: []string{},
	}
	_, err = parseCaddyfile("/etc/caddy/Caddyfile", config)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"status": "error", "message": "Caddyfile parsing failed: %v"}`, err)))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"status": "success", "message": "Configuration reloaded successfully"}`))
}

func iconHandler(w http.ResponseWriter, r *http.Request) {
	// Serve icons from /tmp/icons directory
	filename := filepath.Base(r.URL.Path[7:]) // Remove "/icons/" prefix
	iconPath := filepath.Join("/tmp/icons", filename)
	
	// Set cache headers
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	
	http.ServeFile(w, r, iconPath)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// Serve static files from the static directory
	fs := http.FileServer(http.Dir("static/"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/icons/", iconHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/static/", staticHandler)
	log.Printf("Dashboard running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
