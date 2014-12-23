package context

const (
	STACK    = "fw/rt2/frame"
	SCOPE    = "fw/rt2/scope"
	MOD      = "fw/rt2/module"
	UNIVERSE = "fw/rt2/ctx"
)

type Domain interface {
	ConnectTo(name string, c ContextAware)
	Discover(name string) ContextAware
	Id(c ContextAware) string
	ContextAware
}

type ContextAware interface {
	Init(d Domain)
	Domain() Domain
	Handle(msg interface{})
}
