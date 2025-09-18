# HomeDash

A lightweight Go-based dashboard that automatically generates a homelab dashboard from your Caddyfile.

## Features

- Parses Caddyfile to extract app names and URLs
- Custom icon support via config.yaml with direct image URLs
- Clean, responsive UI with Tailwind CSS
- Minimal Docker image based on Alpine Linux
- Fallback to 🔗 emoji for apps without custom icons
- No authentication required (internal use)

## Quick Start

1. Create a config.yaml file (optional, for custom icons):
   ```yaml
   apps:
     menu:
       icon_url: "https://cdn-icons-png.flaticon.com/512/1827/1827933.png"
     dozzle:
       icon_url: "https://raw.githubusercontent.com/amir20/dozzle/master/assets/logo.png"
   ```

2. Update the volume paths in `docker-compose.yml` to point to your actual files:
   ```yaml
   volumes:
     - /path/to/your/caddy/Caddyfile:/etc/caddy/Caddyfile:ro
     - /path/to/your/dashboard/config.yaml:/etc/config/config.yaml:ro
   ```

3. Build and run:
   ```bash
   docker-compose up --build -d
   ```

4. Access your dashboard at `http://localhost:8080`

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

### Icon Configuration

The `config.yaml` file allows you to specify custom icons for your apps:

```yaml
apps:
  appname:
    icon_url: "https://example.com/path/to/icon.png"
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
@dozzle host dozzle.nelson.red
reverse_proxy @dozzle 192.168.1.100:9999
```

## Icons

Default icons are provided for common apps:
- menu: 📋
- dozzle: 📝  
- books: 📚
- default: 🔗

To restart and pick up Caddyfile changes:
```bash
docker-compose restart
```