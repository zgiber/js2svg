package js2svg

import (
	"encoding/json"
	"fmt"
	"io"
)

type container map[string]interface{}

// GenerateDiagram ...
func GenerateDiagram(src io.Reader, objectPath string) (*Diagram, error) {
	container, err := unmarshalSrc(src)
	schema, found := getField(container, objectPath)
	if !found {
		return nil, fmt.Errorf("schema '%s' not found", objectPath)
	}

	// container, err = resolveReferences(container)
	// if err != nil {
	// 	return nil, err
	// }

	// debug
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	return nil, err
}

// func (c container) resolveReferences(path string) (interface{}, error) {
// 	dst := map[string]interface{}{}
// 	src := getUnknown(c, path)

// 	switch t := src.(type) {
// 	case map[string]interface{}, []interface{}:
// 	default:

// 	}

// 	for k, v := range src {
// 		switch v.(type) {
// 		case map[string]interface{}:
// 			dst[k] = strings.Join([]string{path, k}, ".")
// 		default:
// 			dst
// 		}
// 	}
// 	// for {
// 	// 	// there must be a neater way but i'm in a hurry now, and this tool is not about performance
// 	// 	references := 0
// 	// 	for k, v := range src {
// 	// 		switch t := v.(type) {
// 	// 		case []interface{}:
// 	// 			dst[k] =
// 	// 		case map[string]interface{}:
// 	// 		default:
// 	// 			dst[k] = t
// 	// 		}
// 	// 		if k == "$ref" {
// 	// 			references++
// 	// 			path := strings.ReplaceAll(v.(string)[2:], "/", ".")
// 	// 			value, found := getField(container, path)
// 	// 			if !found {
// 	// 				return nil, fmt.Errorf("invalid reference: %s", v)
// 	// 			}
// 	// 			v = value
// 	// 		}
// 	// 		dst[k] = v
// 	// 	}
// 	// 	if references == 0 {
// 	// 		return resolved, nil
// 	// 	}
// 	// }
// }
