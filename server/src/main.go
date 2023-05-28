package main

import (
	"remote-storage/server/src/file_system"
	"remote-storage/server/src/router"
)

var fs file_system.FileSystem

func main() {
	var route = router.NewHttpRouter()
	route.AddHandler("/test", f)
	route.AddHandler("/upload", f)
	route.AddHandler("/download", f)
	route.AddHandler("/rename", Rename)
	route.AddHandler("/move", Move)
	route.AddHandler("/copy", Copy)
	route.AddHandler("/delete", Delete)

	route.AddHandler("/cd", Cd)
	route.AddHandler("/mkdir", MkDir)
	route.AddHandler("/ls", Ls)
	route.AddHandler("/tree", f)
	route.AddHandler("/find", f)
	route.AddHandler("/show", f)
	route.Listen()
}
