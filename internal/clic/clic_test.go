package clic

import (
	"reflect"
	"testing"
)

type primitives struct {
	unexported string
	Ignored    string `clic:"_"`
	String     string `clic:"a"`
	Int        int    `clic:"bb"`
	Float      float64
	Bool       bool
	Uint       uint
	Ptr        *bool
	Interface  interface{}
}

func TestRunWithPrimitives(t *testing.T) {
	ptrFalse := false
	ptrTrue := true
	p := primitives{
		String:    "foo",
		Int:       42,
		Float:     8.2,
		Bool:      false,
		Uint:      123,
		Ptr:       &ptrFalse,
		Interface: "foo",
	}

	for _, test := range []struct {
		args string
		out  string
		err  string
		c    primitives
	}{
		// GET
		{"", "", "args: empty args", p},
		{"notthere", "", "accessor notthere: struct: field not found", p},
		{"unexported", "", "accessor unexported: struct: field not found", p},
		{"Ignored", "", "accessor Ignored: struct: field not found", p},
		{"a", "foo", "", p},
		{"a.b", "", "accessor b: can't match with type string", p},
		{"bb", "42", "", p},
		{"Float", "8.2", "", p},
		{"Bool", "false", "", p},
		{"Uint", "123", "", p},
		{"Ptr", "false", "", p},
		// {"Interface", "foo", "", p},

		// SET
		{"a bar", "", "", primitives{String: "bar", Int: 42, Float: 8.2, Uint: 123, Ptr: &ptrFalse, Interface: "foo"}},
		{"a \"bar\"", "", "", primitives{String: "bar", Int: 42, Float: 8.2, Uint: 123, Ptr: &ptrFalse, Interface: "foo"}},
		{"bb 12", "", "", primitives{String: "foo", Int: 12, Float: 8.2, Uint: 123, Ptr: &ptrFalse}},
		{"bb string", "", "failed parse value: expected int", primitives{}},
		{"Float 1.1", "", "", primitives{String: "foo", Int: 42, Float: 1.1, Uint: 123, Ptr: &ptrFalse, Interface: "foo"}},
		{"Float string", "", "failed parse value: expected float", primitives{}},
		{"Bool true", "", "", primitives{String: "foo", Int: 42, Float: 8.2, Bool: true, Uint: 123, Ptr: &ptrFalse, Interface: "foo"}},
		{"Bool string", "", "failed parse value: expected bool", primitives{}},
		{"Uint 321", "", "", primitives{String: "foo", Int: 42, Float: 8.2, Uint: 321, Ptr: &ptrFalse, Interface: "foo"}},
		{"Uint string", "", "failed parse value: expected uint", primitives{}},
		{"Ptr false", "", "", primitives{String: "foo", Int: 42, Float: 8.2, Uint: 321, Ptr: &ptrTrue, Interface: "foo"}},
		// {"Interface bar", "", "", primitives{String: "foo", Int: 42, Float: 8.2, Uint: 321, Ptr: &ptrTrue, Interface: "bar"}},
	} {
		t.Run(test.args, func(tt *testing.T) {
			c := test.c
			out, err := Run(test.args, &c)
			if test.err == "" && err != nil {
				tt.Fatalf("clic.Run failed: %s", err.Error())
			}
			if test.err != "" && err != nil && err.Error() != test.err {
				tt.Fatalf("clic.Run should fail with: %s, got: %s", test.err, err)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if !reflect.DeepEqual(test.c, c) {
				tt.Fatalf("config should be %v. got %v", test.c, c)
			}
		})
	}
}

type containers struct {
	Array [2]string
	Slice []string
	Map   map[string]string
}

func TestRunWithContainers(t *testing.T) {
	container := containers{
		Array: [2]string{"a", "b"},
		Slice: []string{"a"},
		Map:   map[string]string{"foo": "bar", "a": "b"},
	}

	for _, test := range []struct {
		args string
		out  string
		err  string
		c    containers
	}{
		// GET
		{"Array", "a\nb", "", container},
		{"Array.a", "", "accessor a: array: expected index accessor", container},
		{"Array[1]", "b", "", container},
		{"Array[2]", "", "accessor 2: array: index out of bounds", container},
		{"Array[b]", "", "accessor b: array: index must be an integer", container},
		{"Slice", "a", "", container},
		{"Slice[0]", "a", "", container},
		{"Slice[1]", "", "accessor 1: array: index out of bounds", container},
		{"Map", "a.b\nfoo.bar", "", container},
		{"Map.foo", "bar", "", container},
		{"Map.blub", "", "accessor blub: map: index not found", container},

		// SET
		{"Array[0] c", "", "", containers{
			Array: [2]string{"c", "b"},
			Slice: []string{"a"},
			Map:   map[string]string{"foo": "bar", "a": "b"},
		}},
		{"Slice[0] b", "", "", containers{
			Array: [2]string{"a", "b"},
			Slice: []string{"b"},
			Map:   map[string]string{"foo": "bar", "a": "b"},
		}},
		{"Slice[] b", "", "", containers{
			Array: [2]string{"a", "b"},
			Slice: []string{"a", "b"},
			Map:   map[string]string{"foo": "bar", "a": "b"},
		}},
		{"Map.foo baar", "", "", containers{
			Array: [2]string{"a", "b"},
			Slice: []string{"a"},
			Map:   map[string]string{"foo": "baar", "a": "b"},
		}},
		{"Map.bla foo", "", "", containers{
			Array: [2]string{"a", "b"},
			Slice: []string{"a"},
			Map:   map[string]string{"foo": "bar", "a": "b", "bla": "foo"},
		}},
	} {
		t.Run(test.args, func(tt *testing.T) {
			c := test.c
			out, err := Run(test.args, &c)
			if test.err == "" && err != nil {
				tt.Fatalf("clic.Run failed: %s", err.Error())
			}
			if test.err != "" && err != nil && err.Error() != test.err {
				tt.Fatalf("clic.Run should fail with: %s, got: %s", test.err, err)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if !reflect.DeepEqual(test.c, c) {
				tt.Fatalf("config should be %v. got %v", test.c, c)
			}
		})
	}
}

type Embeded struct {
	A string
	B string
}

type structs struct {
	Inner primitives
	Embeded
}

func TestRunWithStructs(t *testing.T) {
	st := structs{
		Inner:   primitives{String: "b"},
		Embeded: Embeded{A: "c"},
	}

	for _, test := range []struct {
		args string
		out  string
		err  string
		c    structs
	}{
		// GET
		{"Inner[a]", "", "accessor a: struct: expected named accessor", st},
		{"Inner.a", "b", "", st},
		{"A", "c", "", st},

		// SET
		{"Inner.a c", "", "", structs{
			Inner:   primitives{String: "c"},
			Embeded: Embeded{A: "b"},
		}},
		{"A", "b", "", structs{
			Inner:   primitives{String: "b"},
			Embeded: Embeded{A: "b"},
		}},
	} {
		t.Run(test.args, func(tt *testing.T) {
			c := test.c
			out, err := Run(test.args, &c)
			if test.err == "" && err != nil {
				tt.Fatalf("clic.Run failed: %s", err.Error())
			}
			if test.err != "" && err != nil && err.Error() != test.err {
				tt.Fatalf("clic.Run should fail with: %s, got: %s", test.err, err)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if test.out != out {
				tt.Fatalf("out should be: %s, got: %s", test.out, out)
			}
			if !reflect.DeepEqual(test.c, c) {
				tt.Fatalf("config should be %v. got %v", test.c, c)
			}
		})
	}
}

type argsconfig struct {
	A          string `clic:"a"`
	B          string
	unexported string
	// Pointer    *string
	Array      [3]string
	Slice      []string
	Map        map[string]string
	MapPointer map[string]*simple
	Struct     struct {
		A string
	}
}

type simple struct {
	A string
}

func TestArgs(t *testing.T) {
	c := argsconfig{
		Slice:      []string{"h", "j", "k", "l"},
		Map:        map[string]string{"foo": "", "bar": ""},
		MapPointer: map[string]*simple{"si": &simple{}},
	}

	want := []string{
		"Array",
		"Array[0]",
		"Array[1]",
		"Array[2]",
		"B",
		"Map",
		"Map.bar",
		"Map.foo",
		"MapPointer",
		"MapPointer.si.A",
		"Slice",
		"Slice[0]",
		"Slice[1]",
		"Slice[2]",
		"Slice[3]",
		"Slice[]",
		"Struct.A",
		"a",
	}
	got := Args(c)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got:\n%v\nwant:\n%v", got, want)
	}
}

func TestParseArgs(t *testing.T) {
	for _, test := range []struct {
		args   string
		parsed parsedArgs
		err    string
	}{
		{"", parsedArgs{}, "empty args"},
		{"a", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}}}, ""},
		{"abcdef", parsedArgs{accessors: []accessor{accessor{value: "abcdef", typ: accessorName}}}, ""},
		{"a.", parsedArgs{}, "unexpected trailing dot"},
		{".a", parsedArgs{}, "unexpected leading dot"},
		{"a.b", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}, accessor{value: "b", typ: accessorName}}}, ""},
		{"a 123", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}}}, ""},
		{"a \"foo\"", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}}, value: "foo"}, ""},
		{"a \"bar\"", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}}, value: "bar"}, ""},
		{"a \"foo", parsedArgs{}, "missing closing \""},
		{"a[12]", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}, accessor{value: "12", typ: accessorIndex}}}, ""},
		{"a[12].b", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}, accessor{value: "12", typ: accessorIndex}, accessor{value: "b", typ: accessorName}}}, ""},
		{"a[]", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}, accessor{typ: accessorIndex}}}, ""},
		{"a[12", parsedArgs{}, "missing closing ]"},
		{"a[12] foo", parsedArgs{accessors: []accessor{accessor{value: "a", typ: accessorName}, accessor{value: "12", typ: accessorIndex}}, value: "foo"}, ""},
	} {
		t.Run(test.args, func(tt *testing.T) {
			parsed, err := parseArgs(test.args)
			if test.err == "" && err != nil {
				tt.Fatalf("parseArgs failed: %s", err.Error())
			}
			if test.err != "" && (err == nil || err.Error() != test.err) {
				tt.Fatalf("parseArgs should fail with: %s", test.err)
			}
			if !reflect.DeepEqual(test.parsed.accessors, parsed.accessors) {
				tt.Fatalf("test.parsed.accessors != parsed.accessors: %v != %v", test.parsed.accessors, parsed.accessors)
			}
			if test.parsed.value != test.parsed.value {
				tt.Fatalf("test.parsed.value != parsed.value: %s != %s", test.parsed.value, parsed.value)
			}
		})
	}
}
