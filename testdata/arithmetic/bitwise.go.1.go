// +build example-main

package main

import "fmt"

func main() {
	i := 100
	j := 200

	a := i & j
	a = i & j
	a = i ^ j
	a = i &^ j
	a = i << 1
	a = i >> 1

	fmt.Println(a)
}
