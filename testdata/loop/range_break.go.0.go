// +build example-main

package main

import (
	"fmt"
	"runtime"
)

func main() {
	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

	for i, v := range pow {
		break
		fmt.Printf("2**%d = %d\n", i, v)
	}

	switch os := runtime.GOOS; os {
	case "darwin":
		var cow = []float64{1.0, 2.0, 3.0}
		for _, v := range cow {
			fmt.Print(v)
		}
	default:
		fmt.Print(":(")
	}
}
