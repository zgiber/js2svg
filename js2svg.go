package js2svg

var (
	gapWidth  = 6.0
	gapHeight = 1.0
)

// Object represents a class box in the rendering
type Object struct {
	Name        string
	Description string
	Properties  []Property
	ComposedOf  []Composition
	Position    Position // calculated ... of the top left corner of the rendered object
}

// Property (may rename to field)
type Property struct {
	Name         string
	Description  string
	Relationship string // "0..1" | "1..1" | "1..*"
}

// Composition represents a connection between two class boxes
type Composition struct {
	Relationship string // "0..1" | "1..1" | "1..*"
	Object       *Object
}

// Position is pretty self explanatory
type Position struct {
	X float64
	Y float64
}

// this is called once on the root before the diagram is rendered
func (o *Object) calculateChildPositions() {
	if len(o.ComposedOf) == 0 {
		return
	}

	childPosX := o.Position.X + o.Width() + gapWidth
	for i := 0; i < len(o.ComposedOf); i++ {
		var posY float64
		if i == 0 {
			posY = o.Position.Y
		} else {
			prev := o.ComposedOf[i-1]
			posY = prev.Object.totalHeight() + prev.Object.Position.Y
		}

		o.ComposedOf[i].Object.Position = Position{
			childPosX,
			posY,
		}
		o.ComposedOf[i].Object.calculateChildPositions()
	}
}

// NamePosition returns the postiion where the class name is to be rendered
func (o *Object) NamePosition() Position {
	return Position{
		o.Position.X + o.Width()/2,
		o.Position.Y + 1.3,
	}
}

// FieldPosition is a function used in the template to calculate the positions
// of individual property fields in the class box
func (o *Object) FieldPosition(n int) Position {
	return Position{
		o.Position.X + 1.0,
		o.NamePosition().Y + 2.0 + float64(n)*1.3,
	}
}

// ConnectorInPosition is the target position of an incoming connection (arrow)
func (o *Object) ConnectorInPosition() Position {
	return Position{
		o.Position.X,
		o.Position.Y + 0.5,
	}
}

// ConnectorOutPosition is the source position of an outgoing connection (arrow)
func (o *Object) ConnectorOutPosition() Position {
	return Position{
		o.Position.X + o.Width(),
		o.Position.Y + 0.5,
	}
}

// Width of the class box
func (o *Object) Width() float64 {
	w := len(o.Name)
	for _, p := range o.Properties {
		l := len(p.Name) + len(p.Relationship) + 3
		if l > w {
			w = l
		}
	}
	return float64(w) * 0.8 // for some reason boxes are too wide by default
}

// Height of the class box
func (o *Object) Height() float64 {
	return float64(len(o.Properties))*1.3 + 3.0 // name, line between name and properties, properties, frames
}

// total width of the tree at the widest branch with gaps
func (o *Object) totalWidth() float64 {
	w := o.Position.X + o.Width()
	if len(o.ComposedOf) == 0 {
		return w
	}

	for _, c := range o.ComposedOf {
		if childWidth := c.Object.totalWidth(); childWidth > w {
			w = childWidth
		}
	}

	return w + 1.0
}

// total height of the tree calculated from the lowermost object
func (o *Object) totalHeight() float64 {
	// dfs
	if len(o.ComposedOf) == 0 {
		return o.Height() + gapHeight
	}

	sumLeafNodesHeight := 0.0
	for _, childItem := range o.ComposedOf {
		sumLeafNodesHeight += childItem.Object.totalHeight()
	}

	if o.Height() > sumLeafNodesHeight {
		return o.Height() + gapHeight
	}

	return sumLeafNodesHeight
}
