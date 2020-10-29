package js2svg

import (
	"fmt"
	"html/template"
	"io"
)

const (
	objectFillColor = "azure"
	strokeColor     = "darkslategrey"
	nameColor       = "royalblue"
	propertyColor   = "seagreen"
	connectorColor  = strokeColor

	headers = `<?xml version="1.0" standalone="no"?>
<svg xmlns="http://www.w3.org/2000/svg" font-family="monospace" width="%vem" height="%vem">`
	footer = `</svg>`

	defs = `<defs>
    <marker id="Triangle"
      viewBox="0 0 10 10" refX="0" refY="5" 
      markerUnits="strokeWidth"
      markerWidth="15" markerHeight="10"
      orient="auto">
      <path d="M 0 0 L 10 5 L 0 10 z" />
    </marker>
    
     <marker id="Diamond"
      viewBox="0 0 16 10" refX="0" refY="5" 
      markerUnits="strokeWidth"
      markerWidth="20" markerHeight="10"
      orient="auto">
      <path d="M 0 5 L 8 10 L 16 5 L 8 0 z" />
    </marker>   
</defs>`
)

var (
	objectTemplate = fmt.Sprintf(`
<rect x="{{.Position.X}}em" y="{{.Position.Y}}em" width="{{.Width}}em" height="{{.Height}}em" fill="%s" stroke="%s" stroke-width="2"/>
<text style="font-weight:bold" text-anchor="middle" x="{{.NamePosition.X}}em" y="{{.NamePosition.Y}}em" fill="%s">{{.Name}}<title>{{.Description}}</title></text>
{{range $i, $prop := .Properties}}<text x="{{($.FieldPosition $i).X}}em" y="{{($.FieldPosition $i).Y}}em" fill="%s">{{.Name}} {{.Relationship}}<title>{{.Description}}</title></text>
{{end}}`, objectFillColor, strokeColor, strokeColor, propertyColor)

	connectorTemplate = `
<line x1="{{(index . 0).Start.X}}em" y1="{{(index . 0).Start.Y}}em" x2="{{(index . 0).Stop.X}}em" y2="{{(index . 0).Stop.Y}}em" stroke="darkslategrey" marker-start="url(#Diamond)"/>
<line x1="{{(index . 1).Start.X}}em" y1="{{(index . 1).Start.Y}}em" x2="{{(index . 1).Stop.X}}em" y2="{{(index . 1).Stop.Y}}em" stroke="darkslategrey" />
<line x1="{{(index . 2).Start.X}}em" y1="{{(index . 2).Start.Y}}em" x2="{{(index . 2).Stop.X}}em" y2="{{(index . 2).Stop.Y}}em" stroke="darkslategrey" marker-end="url(#Triangle)"/>
<text x="{{textPosition.X}}em" y="{{textPosition.Y}}em">1..1</text>`
)

// Diagram ...
type Diagram struct {
	root *object
}

// Render ...
func (d *Diagram) Render(dst io.Writer) error {
	// recalculate child positions
	d.root.Position.X = 1 // 1em margin
	d.root.Position.Y = 1 //
	d.root.calculateChildPositions()

	// write the header
	h := fmt.Sprintf(headers, d.root.totalWidth(), d.root.totalHeight()+1.0) // add 1em margin on the bottom
	_, err := dst.Write([]byte(h))
	if err != nil {
		return err
	}

	// defs
	_, err = dst.Write([]byte(defs))
	if err != nil {
		return err
	}

	// render the objects
	if err := renderObject(dst, d.root); err != nil {
		return err
	}

	if err := renderChildObjects(dst, d.root); err != nil {
		return err
	}

	// render the connections
	if err := renderConnections(dst, d.root); err != nil {
		return err
	}

	_, err = dst.Write([]byte(footer))
	return err
}

type line struct {
	Start position
	Stop  position
}

func renderConnection(dst io.Writer, from, to *object) error {
	functions := template.FuncMap(map[string]interface{}{
		"textPosition": func() position {
			return position{
				X: to.Position.X - 3.5,
				Y: to.Position.Y + 0.5,
			}
		},
	})

	// segment points
	sp1 := position{from.Position.X + from.Width(), from.Position.Y + 1.0}
	sp2 := position{sp1.X + 2, sp1.Y}
	sp3 := position{sp2.X, to.Position.Y + 1.0}
	sp4 := position{to.Position.X - 0.8, sp3.Y}

	lines := []line{
		{sp1, sp2},
		{sp2, sp3},
		{sp3, sp4},
	}

	tmpl, err := template.New("connectors").Funcs(functions).Parse(connectorTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(dst, lines); err != nil {
		return err
	}
	return nil
}

func renderObject(dst io.Writer, o *object) error {
	tmpl, err := template.New("object").Parse(objectTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(dst, o); err != nil {
		return err
	}
	return nil
}

func renderChildObjects(dst io.Writer, o *object) error {
	for _, c := range o.composedOf {
		if err := renderObject(dst, c.Object); err != nil {
			return err
		}
		if err := renderChildObjects(dst, c.Object); err != nil {
			return err
		}
	}

	return nil
}

func renderConnections(dst io.Writer, o *object) error {
	for _, c := range o.composedOf {
		// this connection
		renderConnection(dst, o, c.Object)
		// child's connections
		renderConnections(dst, c.Object)
	}
	return nil
}
