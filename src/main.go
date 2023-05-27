package main

import (
	"remote-storage/router"
)

func main() {
	var route = router.NewHttpRouter()
	route.AddHandler("/test", f)
	route.Listen()
}
