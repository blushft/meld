package message

type Message interface {
	Type() string
	Value() Value
}

type Request interface {
	Message

	Handler() string
	Method() string
}

type Response interface {
	Message
}

type Publication interface {
	Message
}

type Event interface {
	Message
}

type RPCRequest struct {
	Service string      `json:"service,omitempty"`
	Handler string      `json:"handler,omitempty"`
	Method  string      `json:"method,omitempty"`
	Args    interface{} `json:"args,omitempty"`
}
