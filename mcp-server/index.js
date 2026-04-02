#!/usr/bin/env node
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const API_BASE = process.env.VEC_API_URL || "http://192.168.83.129:8092";

// ── HTTP helper ──────────────────────────────────────────────────────────────

async function api(method, path, body) {
  const url = `${API_BASE}${path}`;
  const opts = {
    method,
    headers: { "Content-Type": "application/json" },
  };
  if (body !== undefined) {
    opts.body = JSON.stringify(body);
  }
  const res = await fetch(url, opts);
  const contentType = res.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    return { status: res.status, data: await res.json() };
  }
  return { status: res.status, data: await res.text() };
}

// ── safeTool wrapper ─────────────────────────────────────────────────────────
// Catches exceptions from the handler and returns them as MCP error content
// instead of crashing the server.

function safeTool(server, name, description, schema, handler) {
  server.tool(name, description, schema, async (params) => {
    try {
      return await handler(params);
    } catch (err) {
      return {
        content: [{ type: "text", text: `Error: ${err.message || err}` }],
        isError: true,
      };
    }
  });
}

// ── Helpers ──────────────────────────────────────────────────────────────────

function ok(data) {
  const text = typeof data === "string" ? data : JSON.stringify(data, null, 2);
  return { content: [{ type: "text", text }] };
}

// ── MCP Server ───────────────────────────────────────────────────────────────

const server = new McpServer({
  name: "vector-editor",
  version: "1.0.0",
});

// ── health_check ─────────────────────────────────────────────────────────────

safeTool(
  server,
  "health_check",
  "Check if the vector-editor API is running and healthy",
  {},
  async () => {
    const r = await api("GET", "/api/v1/health");
    return ok(r.data);
  }
);

// ── create_document ──────────────────────────────────────────────────────────

safeTool(
  server,
  "create_document",
  "Create a new SVG document with specified dimensions",
  {
    name: z.string().describe("Document name"),
    width: z.number().describe("Width in units"),
    height: z.number().describe("Height in units"),
    unit: z.enum(["px", "mm", "in"]).optional().describe("Unit (default: px)"),
    background: z.string().optional().describe("Background color hex (default: #ffffff)"),
  },
  async ({ name, width, height, unit, background }) => {
    const body = { name, width, height };
    if (unit) body.unit = unit;
    if (background) body.background = background;
    const r = await api("POST", "/api/v1/documents", body);
    return ok(r.data);
  }
);

// ── list_documents ───────────────────────────────────────────────────────────

safeTool(
  server,
  "list_documents",
  "List all SVG documents",
  {},
  async () => {
    const r = await api("GET", "/api/v1/documents");
    return ok(r.data);
  }
);

// ── get_document ─────────────────────────────────────────────────────────────

safeTool(
  server,
  "get_document",
  "Get a document by ID with all its layers",
  {
    document_id: z.string().describe("Document ID"),
  },
  async ({ document_id }) => {
    const r = await api("GET", `/api/v1/documents/${document_id}`);
    return ok(r.data);
  }
);

// ── add_layer ────────────────────────────────────────────────────────────────

safeTool(
  server,
  "add_layer",
  "Add a new layer/shape to a document. Use data map for shape-specific properties (e.g. x, y, width, height for rect; cx, cy, r for circle; d for path; text, font_size for text; points for polygon; x1, y1, x2, y2 for line).",
  {
    document_id: z.string().describe("Document ID"),
    name: z.string().describe("Layer name"),
    type: z
      .enum(["rect", "circle", "ellipse", "path", "text", "polygon", "line", "group", "image"])
      .describe("Shape type"),
    data: z.record(z.string(), z.string()).optional().describe("Shape-specific properties as key-value pairs"),
    fill: z
      .object({
        type: z.enum(["solid", "gradient", "none"]),
        color: z.string().optional(),
        opacity: z.number().optional(),
      })
      .optional()
      .describe("Fill settings"),
    stroke: z
      .object({
        color: z.string().optional(),
        width: z.number().optional(),
        line_cap: z.string().optional(),
        line_join: z.string().optional(),
        dash_array: z.string().optional(),
        opacity: z.number().optional(),
      })
      .optional()
      .describe("Stroke settings"),
    opacity: z.number().optional().describe("Layer opacity 0-1"),
  },
  async ({ document_id, name, type: layerType, data, fill, stroke, opacity }) => {
    const body = { name, type: layerType };
    if (data) body.data = data;
    if (fill) body.fill = fill;
    if (stroke) body.stroke = stroke;
    if (opacity !== undefined) body.opacity = opacity;
    const r = await api("POST", `/api/v1/documents/${document_id}/layers`, body);
    return ok(r.data);
  }
);

// ── update_layer ─────────────────────────────────────────────────────────────

safeTool(
  server,
  "update_layer",
  "Update an existing layer's properties",
  {
    document_id: z.string().describe("Document ID"),
    layer_id: z.string().describe("Layer ID"),
    name: z.string().optional().describe("New name"),
    visible: z.boolean().optional().describe("Visibility"),
    locked: z.boolean().optional().describe("Locked state"),
    opacity: z.number().optional().describe("Opacity 0-1"),
    data: z.record(z.string(), z.string()).optional().describe("Updated shape data"),
  },
  async ({ document_id, layer_id, name, visible, locked, opacity, data }) => {
    const body = {};
    if (name !== undefined) body.name = name;
    if (visible !== undefined) body.visible = visible;
    if (locked !== undefined) body.locked = locked;
    if (opacity !== undefined) body.opacity = opacity;
    if (data) body.data = data;
    const r = await api("PATCH", `/api/v1/documents/${document_id}/layers/${layer_id}`, body);
    return ok(r.data);
  }
);

// ── delete_layer ─────────────────────────────────────────────────────────────

safeTool(
  server,
  "delete_layer",
  "Delete a layer from a document",
  {
    document_id: z.string().describe("Document ID"),
    layer_id: z.string().describe("Layer ID"),
  },
  async ({ document_id, layer_id }) => {
    const r = await api("DELETE", `/api/v1/documents/${document_id}/layers/${layer_id}`);
    return ok(r.data);
  }
);

// ── set_fill ─────────────────────────────────────────────────────────────────

safeTool(
  server,
  "set_fill",
  "Set the fill of a layer",
  {
    document_id: z.string().describe("Document ID"),
    layer_id: z.string().describe("Layer ID"),
    type: z.enum(["solid", "gradient", "none"]).describe("Fill type"),
    color: z.string().optional().describe("Fill color hex"),
    opacity: z.number().optional().describe("Fill opacity 0-1"),
  },
  async ({ document_id, layer_id, type: fillType, color, opacity }) => {
    const fill = { type: fillType };
    if (color) fill.color = color;
    if (opacity !== undefined) fill.opacity = opacity;
    const r = await api("PATCH", `/api/v1/documents/${document_id}/layers/${layer_id}`, { fill });
    return ok(r.data);
  }
);

// ── set_stroke ───────────────────────────────────────────────────────────────

safeTool(
  server,
  "set_stroke",
  "Set the stroke of a layer",
  {
    document_id: z.string().describe("Document ID"),
    layer_id: z.string().describe("Layer ID"),
    color: z.string().optional().describe("Stroke color hex"),
    width: z.number().optional().describe("Stroke width"),
    line_cap: z.string().optional().describe("Line cap: butt, round, square"),
    line_join: z.string().optional().describe("Line join: miter, round, bevel"),
    dash_array: z.string().optional().describe("Dash pattern e.g. '5,3'"),
    opacity: z.number().optional().describe("Stroke opacity 0-1"),
  },
  async ({ document_id, layer_id, color, width, line_cap, line_join, dash_array, opacity }) => {
    const stroke = {};
    if (color) stroke.color = color;
    if (width !== undefined) stroke.width = width;
    if (line_cap) stroke.line_cap = line_cap;
    if (line_join) stroke.line_join = line_join;
    if (dash_array) stroke.dash_array = dash_array;
    if (opacity !== undefined) stroke.opacity = opacity;
    const r = await api("PATCH", `/api/v1/documents/${document_id}/layers/${layer_id}`, { stroke });
    return ok(r.data);
  }
);

// ── set_transform ────────────────────────────────────────────────────────────

safeTool(
  server,
  "set_transform",
  "Set the transform of a layer (translate, rotate, scale, skew)",
  {
    document_id: z.string().describe("Document ID"),
    layer_id: z.string().describe("Layer ID"),
    translate: z.array(z.number()).length(2).optional().describe("[tx, ty]"),
    rotate: z.number().optional().describe("Rotation in degrees"),
    scale: z.array(z.number()).length(2).optional().describe("[sx, sy]"),
    skew_x: z.number().optional().describe("SkewX degrees"),
    skew_y: z.number().optional().describe("SkewY degrees"),
  },
  async ({ document_id, layer_id, translate, rotate, scale, skew_x, skew_y }) => {
    const transform = {};
    if (translate) transform.translate = translate;
    if (rotate !== undefined) transform.rotate = rotate;
    if (scale) transform.scale = scale;
    if (skew_x !== undefined) transform.skew_x = skew_x;
    if (skew_y !== undefined) transform.skew_y = skew_y;
    const r = await api("PATCH", `/api/v1/documents/${document_id}/layers/${layer_id}`, { transform });
    return ok(r.data);
  }
);

// ── render ───────────────────────────────────────────────────────────────────

safeTool(
  server,
  "render",
  "Start a render job for a document. Profiles: png_hires (300dpi), png_web (72dpi), pdf_print (300dpi), svg_source",
  {
    document_id: z.string().describe("Document ID"),
    profile: z
      .enum(["png_hires", "png_web", "pdf_print", "svg_source"])
      .describe("Render profile"),
  },
  async ({ document_id, profile }) => {
    const r = await api("POST", `/api/v1/documents/${document_id}/renders`, { profile });
    return ok(r.data);
  }
);

// ── get_render_status ────────────────────────────────────────────────────────

safeTool(
  server,
  "get_render_status",
  "Check the status and progress of a render job",
  {
    document_id: z.string().describe("Document ID"),
    render_id: z.string().describe("Render job ID"),
  },
  async ({ document_id, render_id }) => {
    const r = await api("GET", `/api/v1/documents/${document_id}/renders/${render_id}`);
    return ok(r.data);
  }
);

// ── get_render_result ────────────────────────────────────────────────────────

safeTool(
  server,
  "get_render_result",
  "Get download information for a completed render. Returns the download URL.",
  {
    document_id: z.string().describe("Document ID"),
    render_id: z.string().describe("Render job ID"),
  },
  async ({ document_id, render_id }) => {
    const statusRes = await api("GET", `/api/v1/documents/${document_id}/renders/${render_id}`);
    const job = statusRes.data;
    if (job.status !== "complete") {
      return ok(`Render not complete yet. Status: ${job.status}, Progress: ${job.progress}%`);
    }
    const downloadUrl = `${API_BASE}/api/v1/documents/${document_id}/renders/${render_id}/download`;
    return ok({ ...job, download_url: downloadUrl });
  }
);

// ── export_svg ───────────────────────────────────────────────────────────────

safeTool(
  server,
  "export_svg",
  "Export a document as raw SVG markup",
  {
    document_id: z.string().describe("Document ID"),
  },
  async ({ document_id }) => {
    const r = await api("POST", `/api/v1/documents/${document_id}/export-svg`);
    return ok(typeof r.data === "string" ? r.data : JSON.stringify(r.data));
  }
);

// ── text_to_path ─────────────────────────────────────────────────────────────

safeTool(
  server,
  "text_to_path",
  "Convert all text elements in a document to vector paths using Inkscape",
  {
    document_id: z.string().describe("Document ID"),
  },
  async ({ document_id }) => {
    const r = await api("POST", `/api/v1/documents/${document_id}/text-to-path`);
    return ok(typeof r.data === "string" ? `Text-to-path complete. SVG size: ${r.data.length} bytes` : r.data);
  }
);

// ── boolean_operation ────────────────────────────────────────────────────────

safeTool(
  server,
  "boolean_operation",
  "Perform a boolean path operation on all shapes in the document (union, intersection, difference, exclusion)",
  {
    document_id: z.string().describe("Document ID"),
    operation: z.enum(["union", "intersection", "difference", "exclusion"]).describe("Boolean operation"),
  },
  async ({ document_id, operation }) => {
    const r = await api("POST", `/api/v1/documents/${document_id}/boolean`, { operation });
    return ok(typeof r.data === "string" ? `Boolean ${operation} complete. SVG size: ${r.data.length} bytes` : r.data);
  }
);

// ── add_text (convenience) ───────────────────────────────────────────────────

safeTool(
  server,
  "add_text",
  "Convenience tool: add a text element to a document with common options",
  {
    document_id: z.string().describe("Document ID"),
    text: z.string().describe("Text content"),
    x: z.string().optional().describe("X position (default: '0')"),
    y: z.string().optional().describe("Y position (default: '0')"),
    font_size: z.string().optional().describe("Font size (default: '16')"),
    font_family: z.string().optional().describe("Font family (default: 'sans-serif')"),
    color: z.string().optional().describe("Text color hex (default: '#000000')"),
    font_weight: z.string().optional().describe("Font weight (default: 'normal')"),
    text_anchor: z.string().optional().describe("Text anchor: start, middle, end"),
  },
  async ({ document_id, text, x, y, font_size, font_family, color, font_weight, text_anchor }) => {
    const data = {
      text,
      x: x || "0",
      y: y || "0",
      font_size: font_size || "16",
      font_family: font_family || "sans-serif",
    };
    if (font_weight) data.font_weight = font_weight;
    if (text_anchor) data.text_anchor = text_anchor;

    const body = {
      name: `text-${text.substring(0, 20)}`,
      type: "text",
      data,
      fill: {
        type: "solid",
        color: color || "#000000",
        opacity: 1,
      },
    };
    const r = await api("POST", `/api/v1/documents/${document_id}/layers`, body);
    return ok(r.data);
  }
);

// ── Start ────────────────────────────────────────────────────────────────────

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("vector-editor MCP server running on stdio");
}

main().catch((err) => {
  console.error("Fatal:", err);
  process.exit(1);
});
