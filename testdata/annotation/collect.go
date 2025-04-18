//go:build examplemain
// +build examplemain

package main

import "fmt"

// mutator-disable-func
func myFunction(a int) int {
	return a * a
}

func regexFunction() {
	fmt.Println("This is a test")
}

func lineFunction() {
	// mutator-disable-next-line numbers/incrementer
	var y = 10

	fmt.Println(y)
	fmt.Println("Line annotation")
}

// mutator-disable-regexp fmt\.Println *
