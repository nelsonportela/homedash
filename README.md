# HomeDash

A lightweight Go-based dashboard that automatically generates a homelab dashboard from your Caddyfile.

## Features

- Parses Caddyfile to extract app names and URLs automatically
- Manual URL entries for external services not in Caddyfile
- Custom icon support via config.yaml with direct image URLs
- Clean, responsive UI with Tailwind CSS
- Minimal Docker image based on Alpine Linux
- Fallback to 🔗 emoji for apps without custom icons
- No authentication required (internal use)

## Docker Compose Files

This project includes two Docker Compose configurations:

- **`docker-compose.yml`** - Production deployment using published Docker image
- **`docker-compose.dev.yml`** - Development setup with local build

## Quick Start

1. Create a config.yaml file (optional, for custom configuration):
   ```yaml
   title: "My HomeLab Dashboard"
   default_theme: "dark"
   
   apps:
     # Caddyfile-based apps (customization only)
     menu:
       icon_url: "https://cdn-icons-png.flaticon.com/512/1827/1827933.png"
       title: "Menu Manager"
       group: "Network"
     dozzle:
       icon_url: "https://raw.githubusercontent.com/amir20/dozzle/master/assets/logo.png"
       title: "Docker Logs"
       group: "Monitoring"
     
     # Manual URL entries (external services)
     adguard:
       url: "https://adguard.vps.example.com"
       icon_url: "https://cdn.jsdelivr.net/gh/AdguardTeam/AdGuardHome@master/internal/home/web/favicon.png"
       title: "AdGuard Home"
       group: "Network"
   ```

2. Update the volume paths in `docker-compose.yml` to point to your actual files:
   ```yaml
   volumes:
     - /path/to/your/caddy/Caddyfile:/etc/caddy/Caddyfile:ro
     - /path/to/your/homedash/config.yaml:/etc/config/config.yaml:ro
   ```

3. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. Access your dashboard at `http://localhost:8080`

## Development

For development with local builds, use the development compose file:

```bash
# Development with local build
docker-compose -f docker-compose.dev.yml up --build -d
```

## Manual Docker Build

```bash
# Build the image
docker build -t homedash .

# Run the container
docker run -d \
  --name homedash \
  -p 8080:8080 \
  -v /path/to/your/caddy/Caddyfile:/etc/caddy/Caddyfile:ro \
  homedash
```

## Configuration

- **PORT**: Environment variable to change the listening port (default: 8080)
- **Caddyfile**: Must be mounted at `/etc/caddy/Caddyfile`
- **Config.yaml**: Optional config file mounted at `/etc/config/config.yaml` for custom icons

### App Configuration

The `config.yaml` file allows you to customize apps and dashboard settings:

```yaml
# Dashboard settings
title: "My HomeLab"
default_theme: "dark"  # or "light"

# App configurations
apps:
  appname:
    icon_url: "https://example.com/path/to/icon.png"
    title: "Custom Display Name"
    group: "Category Name"
    show: true  # Optional: show/hide app (defaults to true)
  hiddenapp:
    title: "Maintenance Tool"
    group: "Development"
    show: false  # This app will be hidden from the dashboard
```

**Available app attributes:**
- `url`: Manual URL for standalone apps not in Caddyfile (required for manual entries)
- `icon_url`: Direct URL to app icon image
- `title`: Custom display name (overrides default from Caddyfile or config key)
- `group`: Category for organizing apps (e.g., "Media", "Monitoring", "Development")
- `show`: Boolean to control app visibility (defaults to `true` if not specified)

### Manual URL Entries

HomeDash supports two types of app entries:

1. **Caddyfile-based apps**: Automatically discovered from your Caddyfile
   - URLs are extracted automatically from Caddyfile patterns
   - Only need configuration for customization (icon, title, group, etc.)

2. **Manual URL entries**: External services not proxied through Caddyfile
   - Must include `url` field with the complete URL
   - Perfect for external VPS services, cloud applications, or other instances
   - Will be mixed with Caddyfile entries in the same groups

Example manual entries:
```yaml
apps:
  adguard:
    url: "https://adguard.vps.example.com"
    icon_url: "https://cdn.jsdelivr.net/gh/AdguardTeam/AdGuardHome@master/internal/home/web/favicon.png"
    title: "AdGuard Home"
    group: "Network"
  grafana:
    url: "https://grafana.example.com"
    icon_url: "https://grafana.com/static/img/menu/grafana2.svg"
    title: "Monitoring Dashboard"
    group: "Monitoring"
```

- Icons are loaded directly from the specified URLs
- If no custom icon is specified, 🔗 emoji is used as fallback
- Supported image formats: PNG, JPG, SVG, etc.
- Icons are displayed at 64x64 pixels

## Supported Caddyfile Format

The app looks for patterns like:
```
@appname host subdomain.domain.com
reverse_proxy @appname ip:port
```

Example:
```
@dozzle host dozzle.example.com
reverse_proxy @dozzle 192.168.1.100:9999
```

## Icon Support

- Icons are loaded directly from URLs specified in config.yaml
- If no custom icon is specified, 🔗 emoji is used as fallback
- Supported image formats: PNG, JPG, SVG, etc.
- Icons are displayed at 64x64 pixels

## Restarting

To restart and pick up Caddyfile changes:
```bash
docker-compose restart
```