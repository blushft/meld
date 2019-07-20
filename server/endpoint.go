package server

type Endpoint interface {
	Path() string
	Service() string
	Handler() string
	Method() string
}
