package main

import (
	"fmt"

	"github.com/imgk/divert-go"
	_ "github.com/imgk/divert-go/resource"
)

func main() {
	ver, err := divert.GetVersionInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(ver)
}
