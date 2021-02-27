package main

import (
	"fmt"
	"os"
)

func main() {
	result, err := calculate(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}
