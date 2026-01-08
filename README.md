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

   groups:
     Network:
       color: "#3B82F6"
       grid_cols: 6
       apps:
         # Caddyfile-based app (no URL needed)
         menu:
           icon_url: "https://cdn-icons-png.flaticon.com/512/1827/1827933.png"
           title: "Menu Manager"

         # Manual URL entry (external service)
         adguard:
           url: "https://adguard.vps.example.com"
           icon_url: "https://cdn.jsdelivr.net/gh/AdguardTeam/AdGuardHome@master/internal/home/web/favicon.png"
           title: "AdGuard Home"

     Monitoring:
       color: "#F59E0B"
       grid_cols: 6
       apps:
         dozzle:
           icon_url: "https://raw.githubusercontent.com/amir20/dozzle/master/assets/logo.png"
           title: "Docker Logs"
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

The `config.yaml` file organizes configuration by groups. Each group controls its display and contains its apps:

```yaml
title: "My HomeLab"
default_theme: "dark"  # or "light"

groups:
  Media:
    color: "#10B981"   # Optional accent color
    grid_cols: 6        # Columns (2-12) on desktop
    apps:
      jellyfin:
        icon_url: "https://cdn.jsdelivr.net/gh/selfhst/icons/webp/jellyfin.webp"
      radarr:
        icon_url: "https://cdn.jsdelivr.net/gh/selfhst/icons/webp/radarr.webp"

  Network:
    color: "#3B82F6"
    grid_cols: 6
    apps:
      adguard:
        url: "https://adguard.example.com"
        icon_url: "https://cdn.jsdelivr.net/gh/AdguardTeam/AdGuardHome@master/internal/home/web/favicon.png"
        title: "AdGuard Home"
      pihole:
        icon_url: "https://cdn.jsdelivr.net/gh/selfhst/icons/webp/pi-hole.webp"
```

Group fields:
- `color`: Hex color for group accent.
- `grid_cols`: Number of columns (2-12) for group width on desktop.
- `apps`: Map of app definitions (keys become app names).

App fields:
- `url`: Manual URL for standalone apps not in Caddyfile (omit for Caddyfile-discovered apps).
- `icon_url`: Direct URL to app icon image.
- `title`: Custom display name (defaults to capitalized key).
- `show`: Boolean to control app visibility (defaults to true).

### Manual URL Entries

HomeDash supports two types of app entries:

1. **Caddyfile-based apps**: Automatically discovered from your Caddyfile
   - URLs are extracted automatically from Caddyfile patterns
   - Only need configuration for customization (icon, title, etc.)

2. **Manual URL entries**: External services not proxied through Caddyfile
   - Must include `url` field with the complete URL
   - Perfect for external VPS services, cloud applications, or other instances

Example manual entries (nested under groups):
```yaml
groups:
  Network:
    apps:
      adguard:
        url: "https://adguard.vps.example.com"
        icon_url: "https://cdn.jsdelivr.net/gh/AdguardTeam/AdGuardHome@master/internal/home/web/favicon.png"
        title: "AdGuard Home"
  Monitoring:
    apps:
      grafana:
        url: "https://grafana.example.com"
        icon_url: "https://grafana.com/static/img/menu/grafana2.svg"
        title: "Monitoring Dashboard"
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