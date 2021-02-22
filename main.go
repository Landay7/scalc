package main

import (
	"fmt"
	"os"
)

func main() {
	result := calculate(os.Args[1:])
	for _, v := range result {
		fmt.Println(v)
	}
}
