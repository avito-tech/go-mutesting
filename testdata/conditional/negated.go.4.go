// +build example-main

package main

import "fmt"

func main() {
	i := 1
	j := 2

	if i > j {
		fmt.Println("1")
	}
	if i < j {
		fmt.Println("2")
	}
	if i >= j {
		fmt.Println("3")
	}
	if i <= j {
		fmt.Println("4")
	}
	if i != j {
		fmt.Println("5")
	}
	if i != j {
		fmt.Println("6")
	}
	fmt.Println("done")
}
