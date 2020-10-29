package js2svg

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type container map[string]interface{}

// GenerateDiagram ...
func GenerateDiagram(src io.Reader, objectPath string) (*Diagram, error) {
	c, err := unmarshalSrc(src)
	if err != nil {
		return nil, err
	}

	// resolve $ref items and replace them with the actual definitions
	resolved, err := c.resolveReferences(objectPath)
	if err != nil {
		return nil, err
	}

	// debug - seems ok...
	// b, err := json.MarshalIndent(resolved, "", "  ")
	// if err != nil {
	// return nil, err
	// }
	// fmt.Println(string(b))

	// expect the top level item to be an object (vs. array or scalar)
	var ok bool
	c, ok = resolved.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("src is not an object")
	}

	psegs := strings.Split(objectPath, ".")
	root := &Object{Name: psegs[len(psegs)-1]}
	err = parseProperties(c, root)
	if err != nil {
		return nil, err
	}

	return &Diagram{Root: root}, nil
}

type iterItem struct {
	Key   string
	Value interface{}
}

// map[string]interface{} yields random order which can be confusing
type iterable []iterItem

func (it iterable) Len() int           { return len(it) }
func (it iterable) Swap(i, j int)      { it[i], it[j] = it[j], it[i] }
func (it iterable) Less(i, j int) bool { return it[i].Key < it[j].Key }

func (it iterable) get(key string) interface{} {
	for _, item := range it {
		if item.Key == key {
			return item.Value
		}
	}
	return nil
}

func mapToIter(m map[string]interface{}) iterable {
	items := make(iterable, 0, len(m))
	for k, v := range m {
		items = append(items, iterItem{Key: k, Value: v})
	}
	sort.Sort(&items)
	return items
}

// expected to pass the root object with its name & description already populated
func parseProperties(c container, parent *Object) error {
	typ, ok := c["type"].(string)
	if !ok {
		return fmt.Errorf("parsing error: expecting an object with 'type' field")
	}

	if typ != "object" {
		return fmt.Errorf("parsing error: container must be of 'object' type (got '%s')", typ)
	}

	properties := getObject(c, "properties")
	for _, prop := range mapToIter(properties) {
		rel := "0..1"

		cc := prop.Value.(map[string]interface{})
		switch cc["type"] {
		case "object":
			if isRequiredParam(c, prop.Key) {
				rel = "1..1"
			}
			child := &Object{}
			child.Name = prop.Key
			if desc := fmt.Sprint(cc["description"]); len(desc) > 0 {
				child.Description = desc
			}
			composeObject(parent, child, rel)
			err := parseProperties(cc, child)
			if err != nil {
				return err
			}

		case "array":
			if isRequiredParam(c, prop.Key) {
				rel = "1..*"
			} else {
				rel = "0..*"
			}

			child := &Object{}
			child.Name = prop.Key
			if desc := fmt.Sprint(cc["description"]); len(desc) > 0 {
				child.Description = desc
			}
			composeObject(parent, child, rel)
			cc = getObject(cc, "items")
			err := parseProperties(cc, child)
			if err != nil {
				return err
			}

		default: // scalar
			if isRequiredParam(c, prop.Key) {
				rel = "1..1"
			} else {
				rel = "0..1"
			}
			info := strings.Builder{}
			for _, key := range []string{"type", "format", "minLength", "maxLength", "description", "enum", "x-namespaced-enum", "pattern"} {
				switch key {
				case "enum", "x-namespaced-enum":
					if value := cc[key]; value != nil {
						info.WriteString("Values:\n")
						for _, v := range value.([]interface{}) {
							info.WriteString(fmt.Sprintf(" - %s\n", v))
						}
					}
				default:
					if cc[key] != nil {
						value := fmt.Sprint(cc[key])
						if len(value) > 0 {
							info.WriteString(fmt.Sprintf("%s: %s\n", strings.Title(key), value))
						}
					}
				}
			}

			parent.Properties = append(parent.Properties, Property{
				Name:         prop.Key,
				Relationship: rel,
				Description:  info.String(),
			})
		}
	}
	return nil
}

func isRequiredParam(c container, name string) bool {
	req, _ := c["required"].([]interface{})
	if req == nil {
		return false
	}

	for _, v := range req {
		if strings.ToLower(v.(string)) == strings.ToLower(name) {
			return true
		}
	}

	return false
}

func composeObject(parent, child *Object, rel string) {
	parent.ComposedOf = append(parent.ComposedOf, Composition{
		Relationship: rel,
		Object:       child,
	})
}

func (c container) resolveReferences(path string) (interface{}, error) {
	src := getUnknown(c, path)
	var err error

	switch t := src.(type) {
	case map[string]interface{}:
		dst := map[string]interface{}{}
		for k, v := range t {
			if k == "$ref" {
				return c.resolveReferences(strings.ReplaceAll(v.(string)[2:], "/", "."))
			}
			subpath := strings.Join([]string{path, k}, ".")
			dst[k], err = c.resolveReferences(subpath)
			if err != nil {
				return nil, err
			}
		}
		return dst, nil

	case []interface{}:
		dst := make([]interface{}, len(t))
		for i := range t {
			subpath := strings.Join([]string{path, fmt.Sprint(i)}, ".")
			dst[i], err = c.resolveReferences(subpath)
			if err != nil {
				return nil, err
			}
		}
		return dst, nil

	case nil:
		return nil, fmt.Errorf("reference '%s' not found", path)

	default:
		return t, nil
	}
}
