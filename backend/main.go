package main

import (
	"easyinvesting/pkg/server"
	"fmt"
)

func main() {
	fmt.Println("Hello")

	s := server.NewServer(8000)
	s.Start()
}
