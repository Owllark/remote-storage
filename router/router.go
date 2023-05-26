package router

type Router interface {
	Listen()
	AddHandler()
}
