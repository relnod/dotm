package clic_test

import (
	"fmt"

	"github.com/relnod/dotm/internal/clic"
)

type myconfig struct {
	A string `clic:"a"`
}

func ExampleRun() {
	c := myconfig{A: "hello"}

	out, err := clic.Run("a", &c)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	_, err = clic.Run(`a "world"`, &c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.A)

	// Output:
	// hello
	// world
}
