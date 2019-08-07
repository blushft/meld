package server

type Router interface {
	Routes() []Route
	Route(string) Route
	AddRoute(Route)
}

type Route interface {
	Method() string
	Path() string
	Name() string
}

type Request interface {
	Path() string
}
