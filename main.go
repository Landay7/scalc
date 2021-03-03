package main

import (
	"fmt"
	"os"
)

func main() {
	result, _, err := calculate(os.Args[1:], 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}
