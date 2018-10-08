package main

import (
	"fmt"

	"github.com/kevinwubert/fb-messenger-analysis/pkg/server"
)

func main() {
	err := server.Main()
	if err != nil {
		fmt.Printf("server errored: %v\n", err)
	}
}
