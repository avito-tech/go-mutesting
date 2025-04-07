//go:build examplemain
// +build examplemain

package main

import (
	"fmt"
	"log/slog"
	"strings"
)

var xx = 42
var xxx = 43

func main() {
	fmt.Println(xx)
	fmt.Println(xxx)

	s := MyStruct{name: "Go"}
	s.Method()

	if strings.Contains("hello", "he") {
		fmt.Println("contains!")
	}
}

type MyStruct struct {
	name string
}

func (m MyStruct) Method() {
	callMe(3, 5)
	fmt.Println("method:", m.name)
}

func callMe(a int, b int) int {
	test := 8
	fmt.Println(test)
	return a + b
}

type Greeter interface {
	Greet() string
}

type Person struct{}

func (p Person) Greet() string {
	slog.Info("structured log")

	return "Hi!"
}
