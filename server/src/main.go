package main

import (
	"remote-storage/server/router"
)

func main() {
	var route = router.NewHttpRouter()
	route.AddHandler("/test", f)
	route.AddHandler("/upload", f)
	route.AddHandler("/download", f)
	route.AddHandler("/rename", f)
	route.AddHandler("/move", f)
	route.AddHandler("/copy", f)

	route.AddHandler("/cd", f)
	route.AddHandler("/tree", f)
	route.AddHandler("/find", f)
	route.AddHandler("/show", f)
	route.Listen()
}
