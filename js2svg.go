package js2svg

var (
	gapWidth  = 6.0
	gapHeight = 1.0
)

type object struct {
	Name        string
	Description string
	Properties  []property
	composedOf  []composition
	Position    position // ... of the top left corner of the rendered object
}

type property struct {
	Name         string
	Description  string
	Relationship string // "0..1" | "1..1" | "1..*"
}

type composition struct {
	Relationship string // "0..1" | "1..1" | "1..*"
	Object       *object
}

type position struct {
	X float64
	Y float64
}

func (o *object) calculateChildPositions() {
	if len(o.composedOf) == 0 {
		return
	}

	childPosX := o.Position.X + o.Width() + gapWidth
	for i := 0; i < len(o.composedOf); i++ {
		var posY float64
		if i == 0 {
			posY = o.Position.Y
		} else {
			prev := o.composedOf[i-1]
			posY = prev.Object.totalHeight() + prev.Object.Position.Y
		}

		o.composedOf[i].Object.Position = position{
			childPosX,
			posY,
		}
		o.composedOf[i].Object.calculateChildPositions()
	}
}

func (o *object) NamePosition() position {
	return position{
		o.Position.X + o.Width()/2,
		o.Position.Y + 1.3,
	}
}

func (o *object) FieldPosition(n int) position {
	return position{
		o.Position.X + 1.0,
		o.NamePosition().Y + 2.0 + float64(n)*1.1,
	}
}

func (o *object) ConnectorInPosition() position {
	return position{
		o.Position.X,
		o.Position.Y + 0.5,
	}
}

func (o *object) ConnectorOutPosition() position {
	return position{
		o.Position.X + o.Width(),
		o.Position.Y + 0.5,
	}
}

func (o *object) Width() float64 {
	w := len(o.Name)
	for _, p := range o.Properties {
		l := len(p.Name) + len(p.Relationship)
		if l > w {
			w = l
		}
	}
	return float64(w)
}

func (o *object) Height() float64 {
	return float64(len(o.Properties) + 4) // name, line between name and properties, properties, frames
}

// returns the total width of the tree
// width of the widest branch with gaps
func (o *object) totalWidth() float64 {
	w := o.Position.X + o.Width()
	if len(o.composedOf) == 0 {
		return w
	}

	for _, c := range o.composedOf {
		if childWidth := c.Object.totalWidth(); childWidth > w {
			w = childWidth
		}
	}

	return w + 1.0
}

// returns the total height of the tree ...
// sum of the height of all leaf nodes plus gaps
func (o *object) totalHeight() float64 {
	// dfs
	if len(o.composedOf) == 0 {
		return o.Height() + gapHeight
	}

	sumLeafNodesHeight := 0.0
	for _, childItem := range o.composedOf {
		sumLeafNodesHeight += childItem.Object.totalHeight()
	}

	if o.Height() > sumLeafNodesHeight {
		return o.Height() + gapHeight
	}

	return sumLeafNodesHeight
}
