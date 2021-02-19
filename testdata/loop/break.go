// +build example-main

package main

import "fmt"

func main() {
	k := 0

	for i := 0; i < 100; i++ {
		if i%2 == 1 {
			k += i
			continue
		}
	}

	for j := 0; j < 400; j++ {
		if j%2 == 1 {
			k += j
			break
		}
	}

	fmt.Println(k)
}
