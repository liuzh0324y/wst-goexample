package main

import (
	"log"

	"github.com/wisdomdev/wisdom-business-server/wst"
)

func main() {
	log.Println("business server.")
	wstserver := wst.New()
	wstserver.Run()
}