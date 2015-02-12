package context

const (
	STACK    = "fw/rt2/frame"
	VSCOPE   = "fw/rt2/scope"
	MOD      = "fw/rt2/module"
	UNIVERSE = "fw/rt2/ctx"
	HEAP     = "fw/rt2/scope,heap"
	MT       = "fw/rt2/table,flow"
	DIGEST   = "fw/cp"

	RETURN = "RETURN"
	META   = "META"
	KEY    = "eval:key"
)

type Factory interface {
	New() Domain
}

type Domain interface {
	Attach(name string, c ContextAware)
	Discover(name string, opts ...interface{}) ContextAware
	Id(c ContextAware) string
	ContextAware
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
