package js2svg

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectDimensions(t *testing.T) {
	// very rudimentary .. for eyaballing results in an svg viewer

	o := testObject("o", 10)
	o1 := testObject("o1", 10)
	o1_1 := testObject("o1_1", 10)
	o1_2 := testObject("o1_2", 10)
	o1_2_1 := testObject("o1_2_1", 10)
	o1_2_2 := testObject("o1_2_2", 10)
	o2 := testObject("o2", 10)

	o1_2.ComposedOf = []Composition{
		{Object: o1_2_1},
		{Object: o1_2_2},
	}

	o1.ComposedOf = []Composition{
		{Object: o1_1, Relationship: "1..1"},
		{Object: o1_2, Relationship: "1..1"},
	}

	o.ComposedOf = []Composition{
		{Object: o1, Relationship: "1..1"},
		{Object: o2, Relationship: "1..1"},
	}

	d := Diagram{Root: o}
	err := d.Render(os.Stdout)
	assert.NoError(t, err)

}

func TestParsing(t *testing.T) {
	// not a real test either...
	f, err := os.Open("test-example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	d, err := ParseToDiagram(f, "components.schemas.OBStandingOrder6Basic")
	if err != nil {
		t.Fatal(err)
	}

	err = d.Render(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func testObject(name string, h int) *Object {
	o := Object{Name: name, Description: fmt.Sprintf("test object %s", name)}
	for i := 0; i < h-4; i++ {
		o.Properties = append(o.Properties, Property{
			Name:         fmt.Sprintf("someField%v", i),
			Description:  fmt.Sprintf("test property field %v", i),
			Relationship: "1..1"})
	}
	return &o
}
