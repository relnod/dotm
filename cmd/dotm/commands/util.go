package commands

import "fmt"

func printl(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
