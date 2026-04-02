package engine

import (
	"encoding/xml"
	"fmt"
	"strings"

	"vector-editor/internal/model"
)

const svgNS = "http://www.w3.org/2000/svg"
const xlinkNS = "http://www.w3.org/1999/xlink"

// SVGDocument is an XML-serializable SVG root element.
type SVGDocument struct {
	XMLName xml.Name    `xml:"svg"`
	XMLNS   string      `xml:"xmlns,attr"`
	Xlink   string      `xml:"xmlns:xlink,attr,omitempty"`
	Width   string      `xml:"width,attr"`
	Height  string      `xml:"height,attr"`
	ViewBox string      `xml:"viewBox,attr,omitempty"`
	Defs    *SVGDefs    `xml:"defs,omitempty"`
	BgRect  *SVGRect    `xml:"rect,omitempty"`
	Groups  []*SVGGroup `xml:"g"`
}

// SVGDefs holds style and gradient definitions.
type SVGDefs struct {
	XMLName xml.Name `xml:"defs"`
	Content string   `xml:",innerxml"`
}

// SVGGroup wraps a <g> element that can contain child shapes.
type SVGGroup struct {
	XMLName   xml.Name `xml:"g"`
	ID        string   `xml:"id,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Display   string   `xml:"display,attr,omitempty"`
	Children  []any    `xml:",any"`
	InnerXML  string   `xml:",innerxml"`
}

// SVGRect represents a <rect> element.
type SVGRect struct {
	XMLName   xml.Name `xml:"rect"`
	ID        string   `xml:"id,attr,omitempty"`
	X         string   `xml:"x,attr,omitempty"`
	Y         string   `xml:"y,attr,omitempty"`
	Width     string   `xml:"width,attr,omitempty"`
	Height    string   `xml:"height,attr,omitempty"`
	RX        string   `xml:"rx,attr,omitempty"`
	RY        string   `xml:"ry,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGCircle represents a <circle> element.
type SVGCircle struct {
	XMLName   xml.Name `xml:"circle"`
	ID        string   `xml:"id,attr,omitempty"`
	CX        string   `xml:"cx,attr,omitempty"`
	CY        string   `xml:"cy,attr,omitempty"`
	R         string   `xml:"r,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGEllipse represents an <ellipse> element.
type SVGEllipse struct {
	XMLName   xml.Name `xml:"ellipse"`
	ID        string   `xml:"id,attr,omitempty"`
	CX        string   `xml:"cx,attr,omitempty"`
	CY        string   `xml:"cy,attr,omitempty"`
	RX        string   `xml:"rx,attr,omitempty"`
	RY        string   `xml:"ry,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGPath represents a <path> element.
type SVGPath struct {
	XMLName   xml.Name `xml:"path"`
	ID        string   `xml:"id,attr,omitempty"`
	D         string   `xml:"d,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGText represents a <text> element.
type SVGText struct {
	XMLName    xml.Name `xml:"text"`
	ID         string   `xml:"id,attr,omitempty"`
	X          string   `xml:"x,attr,omitempty"`
	Y          string   `xml:"y,attr,omitempty"`
	FontFamily string   `xml:"font-family,attr,omitempty"`
	FontSize   string   `xml:"font-size,attr,omitempty"`
	FontWeight string   `xml:"font-weight,attr,omitempty"`
	Fill       string   `xml:"fill,attr,omitempty"`
	Opacity    string   `xml:"opacity,attr,omitempty"`
	Transform  string   `xml:"transform,attr,omitempty"`
	TextAnchor string   `xml:"text-anchor,attr,omitempty"`
	Content    string   `xml:",chardata"`
}

// SVGPolygon represents a <polygon> element.
type SVGPolygon struct {
	XMLName   xml.Name `xml:"polygon"`
	ID        string   `xml:"id,attr,omitempty"`
	Points    string   `xml:"points,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGLine represents a <line> element.
type SVGLine struct {
	XMLName   xml.Name `xml:"line"`
	ID        string   `xml:"id,attr,omitempty"`
	X1        string   `xml:"x1,attr,omitempty"`
	Y1        string   `xml:"y1,attr,omitempty"`
	X2        string   `xml:"x2,attr,omitempty"`
	Y2        string   `xml:"y2,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	StrokeW   string   `xml:"stroke-width,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGImage represents an <image> element.
type SVGImage struct {
	XMLName   xml.Name `xml:"image"`
	ID        string   `xml:"id,attr,omitempty"`
	X         string   `xml:"x,attr,omitempty"`
	Y         string   `xml:"y,attr,omitempty"`
	Width     string   `xml:"width,attr,omitempty"`
	Height    string   `xml:"height,attr,omitempty"`
	Href      string   `xml:"href,attr,omitempty"`
	Opacity   string   `xml:"opacity,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// GenerateSVG creates an SVG XML byte slice from a Document model.
func GenerateSVG(doc *model.Document) ([]byte, error) {
	unit := string(doc.Unit)
	if unit == "" {
		unit = "px"
	}

	svg := SVGDocument{
		XMLNS:  svgNS,
		Xlink:  xlinkNS,
		Width:  fmt.Sprintf("%g%s", doc.Width, unit),
		Height: fmt.Sprintf("%g%s", doc.Height, unit),
	}

	if doc.ViewBox != nil {
		svg.ViewBox = fmt.Sprintf("%g %g %g %g",
			doc.ViewBox.MinX, doc.ViewBox.MinY,
			doc.ViewBox.Width, doc.ViewBox.Height)
	} else {
		svg.ViewBox = fmt.Sprintf("0 0 %g %g", doc.Width, doc.Height)
	}

	// Background rectangle
	if doc.Background != "" && doc.Background != "none" && doc.Background != "transparent" {
		svg.BgRect = &SVGRect{
			Width:  fmt.Sprintf("%g", doc.Width),
			Height: fmt.Sprintf("%g", doc.Height),
			Fill:   doc.Background,
		}
	}

	// Convert layers to SVG groups
	for _, layer := range doc.Layers {
		g := layerToGroup(layer)
		svg.Groups = append(svg.Groups, g)
	}

	output, err := xml.MarshalIndent(&svg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("svg marshal failed: %w", err)
	}

	return append([]byte(xml.Header), output...), nil
}

func layerToGroup(l *model.Layer) *SVGGroup {
	g := &SVGGroup{
		ID: l.ID,
	}

	if !l.Visible {
		g.Display = "none"
	}
	if l.Opacity > 0 && l.Opacity < 1 {
		g.Opacity = fmt.Sprintf("%g", l.Opacity)
	}
	if l.Transform != nil {
		g.Transform = buildTransform(l.Transform)
	}

	// Build inner XML for the shape element
	innerXML := renderShape(l)

	// Render children recursively
	for _, child := range l.Children {
		childG := layerToGroup(child)
		childBytes, err := xml.MarshalIndent(childG, "    ", "  ")
		if err == nil {
			innerXML += "\n" + string(childBytes)
		}
	}

	g.InnerXML = innerXML
	return g
}

func renderShape(l *model.Layer) string {
	fill := "none"
	stroke := "none"
	strokeW := ""
	opacity := ""

	if l.Fill != nil {
		switch l.Fill.Type {
		case model.FillSolid:
			fill = l.Fill.Color
		case model.FillNone:
			fill = "none"
		}
	}
	if l.Stroke != nil {
		stroke = l.Stroke.Color
		if l.Stroke.Width > 0 {
			strokeW = fmt.Sprintf("%g", l.Stroke.Width)
		}
	}
	if l.Opacity > 0 && l.Opacity < 1 {
		opacity = fmt.Sprintf("%g", l.Opacity)
	}

	d := l.Data
	if d == nil {
		d = map[string]string{}
	}

	switch l.Type {
	case model.LayerRect:
		elem := SVGRect{
			ID: l.ID, X: d["x"], Y: d["y"],
			Width: d["width"], Height: d["height"],
			RX: d["rx"], RY: d["ry"],
			Fill: fill, Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerCircle:
		elem := SVGCircle{
			ID: l.ID, CX: d["cx"], CY: d["cy"], R: d["r"],
			Fill: fill, Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerEllipse:
		elem := SVGEllipse{
			ID: l.ID, CX: d["cx"], CY: d["cy"], RX: d["rx"], RY: d["ry"],
			Fill: fill, Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerPath:
		elem := SVGPath{
			ID: l.ID, D: d["d"],
			Fill: fill, Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerText:
		elem := SVGText{
			ID: l.ID, X: d["x"], Y: d["y"],
			FontFamily: d["font_family"], FontSize: d["font_size"],
			FontWeight: d["font_weight"], TextAnchor: d["text_anchor"],
			Fill: fill, Opacity: opacity, Content: d["text"],
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerPolygon:
		elem := SVGPolygon{
			ID: l.ID, Points: d["points"],
			Fill: fill, Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerLine:
		elem := SVGLine{
			ID: l.ID, X1: d["x1"], Y1: d["y1"], X2: d["x2"], Y2: d["y2"],
			Stroke: stroke, StrokeW: strokeW, Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerImage:
		elem := SVGImage{
			ID: l.ID, X: d["x"], Y: d["y"],
			Width: d["width"], Height: d["height"],
			Href: d["href"], Opacity: opacity,
		}
		b, _ := xml.MarshalIndent(elem, "    ", "  ")
		return string(b)

	case model.LayerGroup:
		// Group is handled by children in layerToGroup
		return ""

	default:
		return ""
	}
}

func buildTransform(t *model.Transform) string {
	var parts []string
	if t.Translate != nil {
		parts = append(parts, fmt.Sprintf("translate(%g,%g)", t.Translate[0], t.Translate[1]))
	}
	if t.Rotate != nil {
		parts = append(parts, fmt.Sprintf("rotate(%g)", *t.Rotate))
	}
	if t.Scale != nil {
		parts = append(parts, fmt.Sprintf("scale(%g,%g)", t.Scale[0], t.Scale[1]))
	}
	if t.SkewX != nil {
		parts = append(parts, fmt.Sprintf("skewX(%g)", *t.SkewX))
	}
	if t.SkewY != nil {
		parts = append(parts, fmt.Sprintf("skewY(%g)", *t.SkewY))
	}
	if t.Matrix != nil {
		m := t.Matrix
		parts = append(parts, fmt.Sprintf("matrix(%g,%g,%g,%g,%g,%g)",
			m[0], m[1], m[2], m[3], m[4], m[5]))
	}
	return strings.Join(parts, " ")
}
