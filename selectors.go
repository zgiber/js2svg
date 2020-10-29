package js2svg

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func unmarshalSrc(src io.Reader) (container, error) {
	dst := map[string]interface{}{}
	srcBody, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(srcBody, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}
func getUnknown(m map[string]interface{}, key string) interface{} {
	fieldValue, found := getField(m, key)
	if !found {
		return nil
	}

	return fieldValue
}

func getString(m map[string]interface{}, key string) string {
	fieldValue, _ := getField(m, key)
	if stringValue, ok := fieldValue.(string); ok {
		return stringValue
	}

	return ""
}

func getNumber(m map[string]interface{}, key string) float64 {
	fieldValue, _ := getField(m, key)
	if number, ok := fieldValue.(float64); ok {
		return number
	}

	return 0.0
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	fieldValue, _ := getField(m, key)
	if slice, ok := fieldValue.([]interface{}); ok {
		return slice
	}

	return []interface{}{}
}

func getObject(m map[string]interface{}, key string) map[string]interface{} {
	fieldValue, _ := getField(m, key)
	if obj, ok := fieldValue.(map[string]interface{}); ok {
		return obj
	}

	return map[string]interface{}{}
}

func getField(v interface{}, key string) (interface{}, bool) {
	segment, key := splitKey(key)
	if segment == "" {
		return v, true
	}

	switch t := v.(type) {
	case []interface{}:
		idx, err := strconv.Atoi(segment)
		if err != nil || idx >= len(t) || idx < 0 {
			return nil, false
		}
		return getField(t[idx], key)

	case map[string]interface{}:
		v, found := t[segment]
		if !found {
			return nil, false
		}
		return getField(v, key)
	}
	return nil, false
}

func splitKey(key string) (string, string) {
	key = strings.Trim(key, ".")
	idx := strings.Index(key, ".")
	if idx == -1 {
		return key, ""
	}
	return key[:idx], key[idx+1:]
}
