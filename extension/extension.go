package extension

type EventType int

const (
	NOOP EventType = iota
	OnError
	OnRequest
	OnResponse
	OnRegister
	OnDeregister
	OnPublish
	OnSubscribe
	OnConnect
	OnDisconnect
	OnWatch
	OnUnwatch
	OnChange
	OnUpdate
	OnDelete
	OnQuery
	OnRoute
	OnSelect
	OnFilter
	OnTimeout
)

type Event interface{}

type Extension interface {
	Hook(Event) Event
	HookAsync(Event) <-chan Event
}
