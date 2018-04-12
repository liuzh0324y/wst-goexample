package main

import (
	"log"

	"github.com/wisdomdev/wisdom-business-server/wst"
)

func main() {

	log.Println("business server.")
	wstserver := wst.New("/home/liuzh/work/workspace/src/github.com/wisdomdev/wisdom-business-server/static")
	wstserver.Run()

	// dispatcher := wst.NewEventDispatcher()
	// listener := wst.NewEventListener(EventListener)
	// dispatcher.AddEventListener("hello", listener)

	// time.Sleep(time.Second * 2)

	// dispatcher.DispatchEvent(wst.NewEvent("hello", nil))
}

// func EventListener(event wst.Event) {
// 	fmt.Println(event.Type, event.Object, event.Target)
// }
