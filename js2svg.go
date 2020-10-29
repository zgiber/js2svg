package js2svg

var (
	gapWidth  = 6.0
	gapHeight = 1.0
)

// Object ...
type Object struct {
	Name        string
	Description string
	Properties  []Property
	ComposedOf  []Composition
	Position    Position // ... of the top left corner of the rendered object
}

// Property (may rename to field)
type Property struct {
	Name         string
	Description  string
	Relationship string // "0..1" | "1..1" | "1..*"
}

// Composition ...
type Composition struct {
	Relationship string // "0..1" | "1..1" | "1..*"
	Object       *Object
}

// Position ...
type Position struct {
	X float64
	Y float64
}

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

func (o *Object) NamePosition() Position {
	return Position{
		o.Position.X + o.Width()/2,
		o.Position.Y + 1.3,
	}
}

func (o *Object) FieldPosition(n int) Position {
	return Position{
		o.Position.X + 1.0,
		o.NamePosition().Y + 2.0 + float64(n)*1.1,
	}
}

func (o *Object) ConnectorInPosition() Position {
	return Position{
		o.Position.X,
		o.Position.Y + 0.5,
	}
}

func (o *Object) ConnectorOutPosition() Position {
	return Position{
		o.Position.X + o.Width(),
		o.Position.Y + 0.5,
	}
}

func (o *Object) Width() float64 {
	w := len(o.Name)
	for _, p := range o.Properties {
		l := len(p.Name) + len(p.Relationship)
		if l > w {
			w = l
		}
	}
	return float64(w)
}

func (o *Object) Height() float64 {
	return float64(len(o.Properties) + 4) // name, line between name and properties, properties, frames
}

// returns the total width of the tree
// width of the widest branch with gaps
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

// returns the total height of the tree ...
// sum of the height of all leaf nodes plus gaps
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
