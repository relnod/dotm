// Package clic implements a command line configuration parser.
//
// The package provides a function Run, that can be used to get/set
// configuration values from a struct via a set of given string args.
// To get a list of possible args for a configuration struct you can call the
// Args function.
package clic

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// tagName is the reflection tag name to be used.
const tagName = "clic"

// Run parses the provided args into a set of accessors and an optional value.
// The accessors are separated by a ".". Depending on wether a value was
// provided, it either sets the value or retrives it. Args could look like the
// following:
//      a.b
//      a.b "foo"
//
// An error can occur when trying to find the reflection value, or when setting
// the string value.
//
// It panics for types, that are not implemented in getValue and setValue.
func Run(args string, v interface{}) (string, error) {
	parsedArgs, err := parseArgs(args)
	if err != nil {
		return "", fmt.Errorf("args: %v", err)
	}
	// Check wether v is an interface and retrieve its underlying value.
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return "", fmt.Errorf("non-pointer")
	}
	val, err := findValue(rv.Elem(), parsedArgs.accessors)
	if err != nil {
		return "", err
	}
	if parsedArgs.value == "" {
		return getValue(val), nil
	}
	if !val.CanAddr() {
		return "", fmt.Errorf("not addressable: %v", val)
	}
	return "", setValue(val, parsedArgs.value)
}

// Args returns all valid args, that can be used for the given
// configuration struct v. All valid args can be passed to the Run
// function. The result is ordered in ascending order.
func Args(v interface{}) []string {
	args := findArgs(reflect.ValueOf(v), "")
	sort.Strings(args)
	return args
}

// findArgs collects all valid args for a given reflect.Value.
func findArgs(v reflect.Value, prefix string) []string {
	args := []string{}
	switch v.Type().Kind() {
	case reflect.Ptr:
		return findArgs(v.Elem(), prefix)
	case reflect.Interface:
		return findArgs(reflect.ValueOf(v), prefix)
	case reflect.Struct:
		if prefix != "" {
			prefix += "."
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			// For exported fields the pkg path is empty.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			if field.PkgPath != "" {
				continue
			}
			name := field.Tag.Get(tagName)
			switch name {
			case "_": // Ignore fields, where the tag value is "_".
				continue
			case "": // When the tag value is empty, take the field name.
				name = field.Name
			}
			args = append(args, findArgs(v.Field(i), prefix+name)...)
		}
	case reflect.Array, reflect.Slice:
		args = append(args, prefix)
		if prefix != "" {
			prefix += "."
		}
		for i := 0; i < v.Len(); i++ {
			args = append(args, findArgs(v.Index(i), prefix+strconv.Itoa(i))...)
		}
	case reflect.Map:
		args = append(args, prefix)
		if prefix != "" {
			prefix += "."
		}
		for _, k := range v.MapKeys() {
			args = append(args, findArgs(v.MapIndex(k), prefix+getValue(k))...)
		}
	default:
		args = append(args, prefix)
	}
	return args
}

// findValue tries to find the reflect.Value, that match with the given
// accessors.
// It works by recursively trying to match the current accessor until all
// accessors were matched.
// It returns an error if an accessor cannot be matched.
func findValue(v reflect.Value, accessors []string) (reflect.Value, error) {
	if len(accessors) == 0 {
		return v, nil
	}
	oldAccessors := accessors
	accessor, accessors := accessors[0], accessors[1:]
	switch v.Type().Kind() {
	case reflect.Ptr:
		return findValue(v.Elem(), oldAccessors)
	case reflect.Interface:
		return findValue(reflect.ValueOf(v), oldAccessors)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			// For exported fields the pkg path is empty.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			if field.PkgPath != "" {
				continue
			}
			switch field.Tag.Get(tagName) {
			case "_": // Ignore fields, where the tag value is "_".
				continue
			case "": // When the tag value is empty, try to match the field name.
				if field.Name == accessor {
					return findValue(v.Field(i), accessors)
				}
			case accessor:
				return findValue(v.Field(i), accessors)
			}

			// When the field is an embedded field, check if it contains the
			// accessor. If not proceed with the next field
			if field.Anonymous {
				vv, err := findValue(v.Field(i), oldAccessors)
				if err == nil {
					return vv, nil
				}
			}
		}

		// No struct field has been found for the current accessor.
		return reflect.Value{}, fmt.Errorf("accessor %s: struct: field not found", accessor)
	case reflect.Array, reflect.Slice:
		i, err := strconv.Atoi(accessor)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("accessor %s: array: index must be an integer", accessor)
		}
		if i >= v.Len() {
			return reflect.Value{}, fmt.Errorf("accessor %s: array: index out of bounds", accessor)
		}
		return findValue(v.Index(i), accessors)
	case reflect.Map:
		accessorValue := reflect.ValueOf(accessor)
		for _, k := range v.MapKeys() {
			if accessorValue.String() == k.String() {
				return findValue(v.MapIndex(k), accessors)
			}
		}
		return reflect.Value{}, fmt.Errorf("accessor %s: map: index not found", accessor)
	}

	// We are at value, that doesn't support further propagating the accessor.
	return reflect.Value{}, fmt.Errorf("accessor %s: can't match with type %s", accessor, v.Type().Kind().String())
}

// getValue returns the value of v as a string.
//
// It panics for types, that are not implemented.
func getValue(v reflect.Value) string {
	switch v.Type().Kind() {
	case reflect.Ptr:
		return getValue(v.Elem())
	case reflect.Interface:
		return getValue(reflect.ValueOf(v))
	case reflect.Array, reflect.Slice:
		var sb strings.Builder
		for i := 0; i < v.Len(); i++ {
			sb.WriteString(getValue(v.Index(i)))
			if i < v.Len()-1 {
				sb.WriteString("\n")
			}
		}
		return sb.String()
	case reflect.Map:
		// When getting the string value of a map, first make sure the all map
		// entries are in a deterministic order. This ensures the result is
		// always the same. The entries are ordered by their map key in
		// ascending order.
		entries := []struct {
			k string
			v string
		}{}
		for _, k := range v.MapKeys() {
			entries = append(entries, struct {
				k string
				v string
			}{
				k: getValue(k),
				v: getValue(v.MapIndex(k)),
			})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].k < entries[j].k })
		var sb strings.Builder
		for i, e := range entries {
			sb.WriteString(e.k)
			sb.WriteString(".")
			sb.WriteString(e.v)
			if i < len(entries)-1 {
				sb.WriteString("\n")
			}
		}
		return sb.String()
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	}
	panic(fmt.Sprintf("getValue: type %s not implemented", v.Type().Kind().String()))
}

// setValue sets the value for v. The string value is interpreted according to
// the type of v.
// It Returns an error if the string value can not be converted to the according
// type.
//
// It panics for types, that are not implemented.
func setValue(v reflect.Value, value string) error {
	switch v.Type().Kind() {
	case reflect.Ptr:
		return setValue(v.Elem(), value)
	case reflect.Interface:
		return setValue(reflect.ValueOf(v), value)
	case reflect.String:
		v.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed parse value: expected bool")
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed parse value: expected int")
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed parse value: expected uint")
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed parse value: expected float")
		}
		v.SetFloat(f)
	default:
		panic(fmt.Sprintf("setValue: type %s not implemented", v.Type().Kind().String()))
	}
	return nil
}

// parsedArgs args represents the arguments passed to the Run function.
// The value is optional and is left empty when not set.
// Examples:
//  a.b
//  a.b "b"
type parsedArgs struct {
	accessors []string
	value     string
}

// parseArgs parses string args into a set of accessors separated by a "." and
// an optional value. The args have the following form:
//   accessor(.accessor)* "value"
func parseArgs(args string) (p parsedArgs, err error) {
	if args == "" {
		return parsedArgs{}, fmt.Errorf("empty args")
	}
	// The value is separated from the accessors by the first blank.
	splitted := strings.Split(args, " ")
	if strings.HasSuffix(splitted[0], ".") {
		return parsedArgs{}, fmt.Errorf("unexpected trailing dot")
	}
	p.accessors = strings.Split(splitted[0], ".")

	// Check if a value was set.
	if len(splitted) > 1 {
		// The value can have blanks. So join all remaining splitted values.
		p.value = strings.Join(splitted[1:], " ")

		// The value might be surrounded by quotes. If they are remove them from
		// the value.
		if strings.HasPrefix(p.value, "\"") {
			if !strings.HasSuffix(p.value, "\"") {
				return parsedArgs{}, fmt.Errorf("missing closing quote")
			}
			p.value = strings.Trim(p.value, "\"")
		}
	}

	return p, nil
}
