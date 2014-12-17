package context

const (
	STACK    = "rt2/frame"
	SCOPE    = "rt2/scope"
	MOD      = "rt2/module"
	UNIVERSE = "rt2/ctx"
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
}
