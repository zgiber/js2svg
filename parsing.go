package js2svg

import (
	"encoding/json"
	"fmt"
	"io"
)

// GenerateDiagram ...
func GenerateDiagram(src io.Reader, objectPath string) (*Diagram, error) {
	container, err := unmarshalSrc(src)
	schema, found := getField(container, objectPath)
	if !found {
		return nil, fmt.Errorf("schema '%s' not found", objectPath)
	}

	// debug
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	return nil, err
}
