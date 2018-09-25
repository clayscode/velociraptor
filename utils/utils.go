package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"www.velocidex.com/golang/vfilter"
)

func InString(hay *[]string, needle string) bool {
	for _, x := range *hay {
		if x == needle {
			return true
		}
	}

	return false
}

func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

// Massage a windows path into a standard form:
// \ are replaced with /
// Drive letters are preceeded with /
// Example: c:\windows ->  /c:/windows
func Normalize_windows_path(filename string) string {
	filename = strings.Replace(filename, "\\", "/", -1)
	if !strings.HasPrefix(filename, "/") {
		filename = "/" + filename
	}
	return filename
}

func hard_wrap(text string, colBreak int) string {
	text = strings.TrimSpace(text)
	text = strings.Replace(text, "\r\n", "\n", -1)
	wrapped := ""
	var i int
	for i = 0; len(text[i:]) > colBreak; i += colBreak {

		wrapped += text[i:i+colBreak] + "\n"

	}
	wrapped += text[i:]

	return wrapped
}

func Stringify(value interface{}, scope *vfilter.Scope) string {
	// Deal with pointers to things as those things.
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return Stringify(reflect.Indirect(
			reflect.ValueOf(value)).Interface(), scope)
	}

	json_marshall := func(value interface{}) string {
		if k, err := json.Marshal(value); err == nil {
			if len(k) > 0 && k[0] == '"' && k[len(k)-1] == '"' {
				k = k[1 : len(k)-1]
			}

			return hard_wrap(string(k), 30)
		}
		return ""
	}

	switch t := value.(type) {
	case vfilter.Dict:
		result := []string{}
		iter := t.IterFunc()
		for kv, ok := iter(); ok; kv, ok = iter() {
			result = append(result, fmt.Sprintf("%v: %v", kv.Key, kv.Value))
		}
		return strings.Join(result, "\n")

	case vfilter.StringProtocol:
		return t.ToString(scope)

	case []byte:
		return hard_wrap(string(t), 30)

	case string:
		return hard_wrap(t, 30)

	case json.Marshaler:
		return json_marshall(value)
	default:
		// For normal structs json is a pretty good encoder.
		if reflect.TypeOf(value).Kind() == reflect.Struct {
			return json_marshall(value)
		}

		// Anything else we output something useful.
		return hard_wrap(fmt.Sprintf("%v", value), 30)
	}
}

func SlicesEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, a_item := range a {
		if a_item != b[idx] {
			return false
		}
	}

	return true
}
