package context

const (
	STACK    = "fw/rt2/frame"
	SCOPE    = "fw/rt2/scope"
	MOD      = "fw/rt2/module"
	UNIVERSE = "fw/rt2/ctx"
	HEAP     = "fw/rt2/scope,heap"
	MT       = "fw/rt2/table,flow"
	DIGEST   = "fw/cp"

	RETURN = "RETURN"
	KEY    = "KEY"
)

type Factory interface {
	New() Domain
}

type Domain interface {
	Attach(name string, c ContextAware)
	Discover(name string) ContextAware
	Id(c ContextAware) string
	ContextAware
	Global() Domain
}

type ContextAware interface {
	Init(d Domain)
	Domain() Domain
	Handle(msg interface{})
}

type data struct {
	inner interface{}
}

func (d *data) Init(x Domain)          {}
func (d *data) Domain() Domain         { return nil }
func (d *data) Handle(msg interface{}) {}

func Data(x interface{}) ContextAware {
	return &data{inner: x}
}
