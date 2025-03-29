//go:build windows && (amd64 || 386 || arm64)

package main

import (
	"fmt"

	"github.com/imgk/divert-go"
)

func main() {
	ver, err := divert.GetVersionInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(ver)
}
