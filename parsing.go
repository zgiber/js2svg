package js2svg

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

var (
	// ExternalDivider can be set to custom character in places where . is escaped
	ExternalDivider = "."
	internalDivider = "."
)

// ParseToDiagram performs all the necessary steps for creating a diagram in one function.
func ParseToDiagram(src io.Reader, objectPath string) (*Diagram, error) {
	objectPath = strings.ReplaceAll(objectPath, ExternalDivider, internalDivider)
	m, err := ParseToMap(src, objectPath)
	if err != nil {
		return nil, err
	}
	return MakeDiagram(m, objectPath)
}

// MakeDiagram from a document unmarshalled to a map (useful when multiple diagrams are rendered from the same document)
func MakeDiagram(m map[string]interface{}, path string) (*Diagram, error) {
	psegs := strings.Split(path, ExternalDivider)
	root := &Object{Name: psegs[len(psegs)-1]}
	err := parseProperties(m, root)
	if err != nil {
		return nil, err
	}

	return &Diagram{Root: root}, nil
}

// ParseToMap the selected objectPath. If ObjectPath is the root element of the jsonschema document
// then all references are going to be resolved in the returned map. This map can be resued with MakeDiagram
// to generate multiple diagrams from the same source without unmarshallig / resolving references each time.
func ParseToMap(src io.Reader, objectPath string) (map[string]interface{}, error) {
	objectPath = strings.ReplaceAll(objectPath, ExternalDivider, internalDivider)
	c, err := unmarshalSrc(src)
	if err != nil {
		return nil, err
	}

	// resolve $ref items and replace them with the actual definitions
	resolved, err := resolveReferences(c, objectPath)
	if err != nil {
		return nil, err
	}

	// expect the top level item to be an object (vs. array or scalar)
	var ok bool
	c, ok = resolved.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("src is not an object")
	}
	return c, nil

}

// map[string]interface{} is iterated in random order which can be confusing.
// Iterable can provide the same values as slice members sorted alphabetically
type iterable []iterItem
type iterItem struct {
	Key   string
	Value interface{}
}

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

// expected to pass a root object with its name & description populated outside this func
func parseProperties(m map[string]interface{}, parent *Object) error {
	typ, ok := m["type"].(string)
	if !ok {
		return fmt.Errorf("parsing error: expecting an object with 'type' field: %v", m)
	}

	if typ != "object" {
		return nil
		// return fmt.Errorf("parsing error: container must be of 'object' type (got '%s')", typ)
	}

	properties := GetObject(m, "properties")
	if properties == nil {
		fmt.Println("not an object")
		return nil
	}
	for _, prop := range mapToIter(properties) {
		rel := "0..1"

		cm := prop.Value.(map[string]interface{})
		switch cm["type"] {
		case "object":
			if isRequiredField(m, prop.Key) {
				rel = "1..1"
			}
			child := &Object{}
			child.Name = prop.Key
			if desc, isSet := cm["description"].(string); isSet && len(desc) > 0 {
				child.Description = desc
			}
			composeObject(parent, child, rel)
			err := parseProperties(cm, child)
			if err != nil {
				return err
			}

		case "array":
			if isRequiredField(m, prop.Key) {
				rel = "1..*"
			} else {
				rel = "0..*"
			}

			child := &Object{}
			child.Name = prop.Key
			if desc, isSet := cm["description"].(string); isSet && len(desc) > 0 {
				child.Description = desc
			}
			cm = GetObject(cm, "items")
			switch cm["type"].(string) {
			case "array", "object":
				composeObject(parent, child, rel)
				// new object within array
				err := parseProperties(cm, child)
				if err != nil {
					return err
				}
				continue
			default:
				// properties := GetObject(cm, "properties")

				setScalarProperty(prop.Key, rel, properties, parent)
			}

		default: // scalar
			if isRequiredField(m, prop.Key) {
				rel = "1..1"
			} else {
				rel = "0..1"
			}
			setScalarProperty(prop.Key, rel, cm, parent)
		}
	}
	return nil
}

func setArrayProperties(itemsSchema map[string]interface{}) {
	// WIP: refactor parseProperties
}

func setObjectProperties(objectSchema map[string]interface{}, object *Object) error {
	// WIP: refactor parseProperties
	m := GetObject(objectSchema, "properties")
	if m == nil {
		return nil
	}

	properties := mapToIter(m)
	_ = properties

	// for each property:
	// if object, create child and do recursion
	// if array, create child and do recursion with "items"
	// in any other case, apply as scalar property

	return nil
}

func setScalarProperty(propertyName, rel string, propertySchema map[string]interface{}, o *Object) {
	info := strings.Builder{}
	for _, key := range []string{"type", "format", "minLength", "maxLength", "description", "enum", "x-namespaced-enum", "pattern"} {
		switch key {
		case "enum", "x-namespaced-enum":
			if value := propertySchema[key]; value != nil {
				info.WriteString("Values:\n")
				for _, v := range value.([]interface{}) {
					info.WriteString(fmt.Sprintf(" - %s\n", v))
				}
			}
		default:
			if propertySchema[key] != nil {
				value := fmt.Sprint(propertySchema[key])
				if len(value) > 0 {
					info.WriteString(fmt.Sprintf("%s: %s\n", strings.Title(key), value))
				}
			}
		}
	}

	o.Properties = append(o.Properties, Property{
		Name:         propertyName,
		Relationship: rel,
	})
	if desc := info.String(); len(desc) > 0 {
		o.Description = desc
	}
}

// m is the map with the unmarshalled root schema of the object with the property
// containing the "required" field. name is the name of the field.
func isRequiredField(m map[string]interface{}, name string) bool {
	req, _ := m["required"].([]interface{})
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

func resolveReferences(c map[string]interface{}, path string) (interface{}, error) {
	src := GetUnknown(c, path)
	var err error

	switch t := src.(type) {
	case map[string]interface{}:
		dst := map[string]interface{}{}
		for k, v := range t {
			if k == "$ref" {
				return resolveReferences(c, strings.ReplaceAll(v.(string)[2:], "/", "."))
			}
			subpath := strings.Join([]string{path, k}, ".")
			dst[k], err = resolveReferences(c, subpath)
			if err != nil {
				return nil, err
			}
		}
		return dst, nil

	case []interface{}:
		dst := make([]interface{}, len(t))
		for i := range t {
			subpath := strings.Join([]string{path, fmt.Sprint(i)}, ".")
			dst[i], err = resolveReferences(c, subpath)
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
