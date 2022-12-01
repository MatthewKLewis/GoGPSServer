package main

import (
	"fmt"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
