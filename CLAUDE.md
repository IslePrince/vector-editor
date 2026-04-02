# Vector Editor MCP

SVG vector graphics editor with an HTTP API and MCP server. Wraps Inkscape CLI for advanced operations (render, text-to-path, boolean ops) and uses pure Go for SVG generation.

## Architecture

- **Go HTTP API** (Chi router) on port 8092 -- document/layer CRUD, render pipeline, SVG export, Inkscape operations
- **MCP Server** (Node.js, `@modelcontextprotocol/sdk`) -- stdio bridge that wraps the HTTP API for Claude
- **Inkscape CLI** -- used for PNG/PDF export, text-to-path conversion, boolean path operations
- **Storage** -- JSON files on disk with in-memory index; documents stored in `data/documents/`
- **Render Worker Pool** -- async job queue for export operations

## Host

The API runs on Linux at **192.168.83.129:8092**. It does NOT run on Windows.

## Project Structure

```
cmd/server/main.go          Entry point
internal/
  api/                       HTTP handlers + Chi router
  config/                    Environment-based configuration (VEC_ prefix)
  model/                     Domain types: Document, Layer, Fill, Stroke, Transform, RenderJob
  engine/inkscape.go         Inkscape CLI wrapper (exec.CommandContext)
  engine/svg.go              Pure Go SVG generation via encoding/xml
  render/                    Render profiles, Inkscape renderer, async worker pool
  storage/                   Storage interface + local JSON file implementation
  queue/                     Job queue interface + in-memory implementation
mcp-server/                  Node.js MCP server (index.js + package.json)
docker/Dockerfile            Multi-stage build (golang:1.22-bookworm -> debian:bookworm-slim)
docker-compose.yml           Docker Compose with ve_data volume
```

## Environment Variables

All prefixed with `VEC_`:

| Variable | Default | Description |
|---|---|---|
| VEC_PORT | 8092 | HTTP port |
| VEC_DATA_DIR | ./data | Data directory for documents and renders |
| VEC_INKSCAPE_PATH | inkscape | Path to Inkscape binary |
| VEC_MAX_WORKERS | 4 | Render worker pool size |

## Docker

```bash
docker compose up --build -d
```

Builds a multi-stage image: Go 1.22 builder, then debian:bookworm-slim runtime with Inkscape and Google Fonts.

## Capabilities

- **Shapes**: rect, circle, ellipse, path, text, polygon, line, group, image
- **Fill**: solid (hex color), gradient, none -- with opacity
- **Stroke**: color, width, line_cap, line_join, dash_array, opacity
- **Transforms**: translate, rotate, scale, skewX, skewY, matrix
- **Render profiles**: png_hires (300dpi), png_web (72dpi), pdf_print (300dpi), svg_source
- **Boolean ops**: union, intersection, difference, exclusion (via Inkscape)
- **Text to path**: converts all text to vector outlines (via Inkscape)

## MCP Tools

health_check, create_document, list_documents, get_document, add_layer, update_layer, delete_layer, set_fill, set_stroke, set_transform, render, get_render_status, get_render_result, export_svg, text_to_path, boolean_operation, add_text

## API Routes

```
GET    /api/v1/health
GET    /api/v1/capabilities
POST   /api/v1/documents
GET    /api/v1/documents
GET    /api/v1/documents/{id}
PATCH  /api/v1/documents/{id}
DELETE /api/v1/documents/{id}
POST   /api/v1/documents/{id}/layers
GET    /api/v1/documents/{id}/layers
GET    /api/v1/documents/{id}/layers/{lid}
PATCH  /api/v1/documents/{id}/layers/{lid}
DELETE /api/v1/documents/{id}/layers/{lid}
POST   /api/v1/documents/{id}/renders
GET    /api/v1/documents/{id}/renders/{rid}
GET    /api/v1/documents/{id}/renders/{rid}/download
POST   /api/v1/documents/{id}/export-svg
POST   /api/v1/documents/{id}/text-to-path
POST   /api/v1/documents/{id}/boolean
```
