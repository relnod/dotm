// Package clic implements a command line configuration parser.
//
// The package provides a function Run, that can be used to get/set
// configuration values from a struct via a set of given string args.
// To get a list of possible args for a configuration struct you can call the
// Args function.
package clic

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// tagName is the reflection tag name to be used.
const tagName = "clic"

// Run parses the provided args into a set of accessors and an optional value.
// The accessors are separated by a ".". Depending on wether a value was
// provided it visitor either sets the value or retrives it. Args could look like the
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
	var visitor valueVisitor
	if parsedArgs.value == "" {
		visitor = &valueGetter{}
	} else {
		visitor = &valueSetter{value: parsedArgs.value}
	}
	err = iterateAccessors(rv.Elem(), parsedArgs.accessors, visitor)
	if err != nil && err != io.EOF {
		return "", err
	}
	return visitor.String(), nil
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
	case reflect.Slice:
		args = append(args, prefix+"[]")
		fallthrough
	case reflect.Array:
		args = append(args, prefix)
		for i := 0; i < v.Len(); i++ {
			args = append(args, findArgs(v.Index(i), prefix+"["+strconv.Itoa(i)+"]")...)
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

type valueVisitor interface {
	Value(v reflect.Value, accessors []accessor) error
	String() string
}

func iterateAccessors(v reflect.Value, accessors []accessor, visitor valueVisitor) error {
	err := visitor.Value(v, accessors)
	if err != nil {
		return err
	}
	if len(accessors) == 0 {
		// Abort the recursion, when all accessorve been processed.
		return nil
	}

	oldAccessors := accessors
	accessor, accessors := accessors[0], accessors[1:]
	switch v.Type().Kind() {
	case reflect.Ptr:
		return iterateAccessors(v.Elem(), oldAccessors, visitor)
	case reflect.Interface:
		return iterateAccessors(reflect.ValueOf(v), oldAccessors, visitor)
	case reflect.Struct:
		if accessor.typ != accessorName {
			return fmt.Errorf("accessor %s: struct: expected named accessor", accessor.value)
		}
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
				if field.Name == accessor.value {
					return iterateAccessors(v.Field(i), accessors, visitor)
				}
			case accessor.value:
				return iterateAccessors(v.Field(i), accessors, visitor)
			}

			// When the field is an embedded field, check if visitor contains the
			// accessor. If not proceed with the next field.
			if field.Anonymous {
				// A recorder is used, to check if the embedded field
				// matches the accessor. This prevents iterations leaking to the
				// main valueVisitor, if visitor doesn't match.
				r := &recorder{records: make([]record, 0), visitor: visitor}
				err := iterateAccessors(v.Field(i), oldAccessors, r)
				if err == nil {
					return r.Playback()
				}
			}
		}

		// No struct field has been found for the current accessor.
		return fmt.Errorf("accessor %s: struct: field not found", accessor.value)
	case reflect.Array, reflect.Slice:
		if accessor.typ != accessorIndex {
			return fmt.Errorf("accessor %s: array: expected index accessor", accessor.value)
		}
		i, err := strconv.Atoi(accessor.value)
		if err != nil {
			return fmt.Errorf("accessor %s: array: index must be an integer", accessor.value)
		}
		if i >= v.Len() {
			return fmt.Errorf("accessor %s: array: index out of bounds", accessor.value)
		}
		return iterateAccessors(v.Index(i), accessors, visitor)
	case reflect.Map:
		accessorValue := reflect.ValueOf(accessor.value)
		for _, k := range v.MapKeys() {
			if accessorValue.String() == k.String() {
				return iterateAccessors(v.MapIndex(k), accessors, visitor)
			}
		}
		return fmt.Errorf("accessor %s: map: index not found", accessor.value)
	}

	// We are at value, that doesn't support further propagating the accessor.
	return fmt.Errorf("accessor %s: can't match with type %s", accessor.value, v.Type().Kind().String())
}

type valueSetter struct {
	value string
}

func (visitor *valueSetter) Value(v reflect.Value, accessors []accessor) error {
	switch len(accessors) {
	case 0:
		// If we are at the last reflect.Value, try to set the value.
		if !v.CanAddr() {
			return fmt.Errorf("not addressable: %s: %v", v.Type().Kind(), v)
		}
		return setValue(v, visitor.value)
	case 1:
		// If we are at the second last accessor, we might set the value
		// already, for some special cased types.
		switch v.Type().Kind() {
		case reflect.Slice:
			if accessors[0].typ == accessorIndex && accessors[0].value == "" {
				val, err := newValue(v.Type().Elem().Kind(), visitor.value)
				if err != nil {
					return err
				}
				reflect.Append(v, val)
				return io.EOF
			}
		case reflect.Map:
			val, err := newValue(v.Type().Elem().Kind(), visitor.value)
			if err != nil {
				return err
			}
			v.SetMapIndex(reflect.ValueOf(accessors[0].value), val)
			return io.EOF
		}
	}
	return nil
}

func (visitor *valueSetter) String() string { return "" }

type valueGetter struct {
	value string
}

func (visitor *valueGetter) Value(v reflect.Value, accessors []accessor) error {
	if len(accessors) == 0 {
		visitor.value = getValue(v)
		return io.EOF
	}
	return nil
}

func (visitor *valueGetter) String() string { return visitor.value }

type record struct {
	v         reflect.Value
	accessors []accessor
}

type recorder struct {
	records []record
	visitor valueVisitor
}

func (visitor *recorder) Value(v reflect.Value, accessors []accessor) error {
	visitor.records = append(visitor.records, record{v: v, accessors: accessors})
	return nil
}

func (visitor *recorder) String() string { return "" }

func (visitor *recorder) Playback() error {
	for _, r := range visitor.records {
		err := visitor.visitor.Value(r.v, r.accessors)
		if err != nil {
			return err
		}
	}
	return nil
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

// newValue creates a new reflect.Value for the given string value.
// Tries to convert the string value according to the given reflect.Kind.
//
// It panics for types, that are not implemented.
func newValue(k reflect.Kind, value string) (reflect.Value, error) {
	switch k {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed parse value: expected bool")
		}
		return reflect.ValueOf(b), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed parse value: expected int")
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed parse value: expected uint")
		}
		return reflect.ValueOf(i), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed parse value: expected float")
		}
		return reflect.ValueOf(f), nil
	}
	panic(fmt.Sprintf("newValue: type %s not implemented", k.String()))
}

type accessorType int

const (
	accessorName accessorType = iota
	accessorIndex
)

type accessor struct {
	value string
	typ   accessorType
}

// parsedArgs represents the parsed arguemnts. It has a list of accessors and an
// optional value.
type parsedArgs struct {
	accessors []accessor
	value     string
}

// parseState defines the state of the argument parser. Additional states can be
// created by combining existing states.
type parseState int

const (
	parseName parseState = 1 << iota
	parseIndex
	parseValue
	parseValueQuoted
	parseEnd
)

// parseArgs parses the given string arguments into a parsedArgs struct.
// The parsing works with a simple state machine with the current character as
// input.
// It returns an error, when the arguments are in an invalid format.
func parseArgs(args string) (p parsedArgs, err error) {
	if args == "" {
		return parsedArgs{}, fmt.Errorf("empty args")
	}

	state := parseName | parseIndex // Set start state to parseName or parseIndex
	tokenStart := 0                 // tokenStart tracks the beginning of a token.
	addAccessor := func(typ accessorType, end int) {
		value := ""
		if tokenStart < end {
			value = args[tokenStart:end]
		}
		p.accessors = append(p.accessors, accessor{
			typ:   typ,
			value: value,
		})
	}

	for i := 0; i < len(args); i++ {
		switch state {
		case parseName | parseIndex | parseValue:
			switch args[i] {
			case ' ':
				state = parseValue | parseValueQuoted
				tokenStart = i + 1
				continue
			}
			fallthrough
		case parseName | parseIndex:
			switch args[i] {
			case '[':
				state = parseIndex
				tokenStart = i + 1
			case '.':
				if i == 0 {
					return parsedArgs{}, fmt.Errorf("unexpected leading dot")
				}
				state = parseName
				tokenStart = i + 1
			default:
				state = parseName
				tokenStart = i
			}
		case parseName:
			switch args[i] {
			case '[':
				addAccessor(accessorName, i)
				state = parseIndex
				tokenStart = i + 1
			case '.':
				if i == len(args)-1 {
					return parsedArgs{}, fmt.Errorf("unexpected trailing dot")
				}
				addAccessor(accessorName, i)
				state = parseName
				tokenStart = i + 1
			case ' ':
				addAccessor(accessorName, i)
				tokenStart = i + 1
				state = parseValue | parseValueQuoted
			}
		case parseIndex:
			switch args[i] {
			case ']':
				addAccessor(accessorIndex, i)
				state = parseName | parseIndex | parseValue
			}
		case parseValue | parseValueQuoted:
			switch args[i] {
			case '"':
				state = parseValueQuoted
				tokenStart = i + 1
			default:
				state = parseValue
			}
		case parseValue: // Do nothing
		case parseValueQuoted:
			switch args[i] {
			case '"':
				if i == len(args)-1 {
					p.value = args[tokenStart:i]
					state = parseEnd
				}
			}
		}
	}

	// After processing all inputs, check the state of the parser and either do
	// nothing, return an error for unfinished tokens or add a new token.
	switch state {
	case parseName:
		addAccessor(accessorName, len(args))
	case parseIndex:
		return parsedArgs{}, fmt.Errorf("missing closing ]")
	case parseValue:
		p.value = args[tokenStart:len(args)]
	case parseValueQuoted:
		return parsedArgs{}, fmt.Errorf("missing closing \"")
	}
	return p, nil
}
