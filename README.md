# HomeDash

A lightweight Go-based dashboard that automatically generates a homelab dashboard from your Caddyfile.

## Features

- Parses Caddyfile to extract app names and URLs
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
     menu:
       icon_url: "https://cdn-icons-png.flaticon.com/512/1827/1827933.png"
       title: "Menu Manager"
       group: "Network"
     dozzle:
       icon_url: "https://raw.githubusercontent.com/amir20/dozzle/master/assets/logo.png"
       title: "Docker Logs"
       group: "Monitoring"
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
```

**Available app attributes:**
- `icon_url`: Direct URL to app icon image
- `title`: Custom display name (overrides default from Caddyfile)
- `group`: Category for organizing apps (e.g., "Media", "Monitoring", "Development")

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