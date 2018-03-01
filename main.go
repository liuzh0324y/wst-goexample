package main

import (
	"fmt"
	"github.com/wisdomdev/wisdom-business-server/wst"
)

func main() {
	fmt.Println("business server.")
	httpserver := wst.New()
	httpserver.Run()
}