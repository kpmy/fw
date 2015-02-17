package context

const (
	UNIVERSE = "ctx"
	STACK    = "stack"
	MOD      = "mt"
	MT       = "flow"
	DIGEST   = "cp"
	HEAP     = "heap"
	SCOPE    = "mod"
	CALL     = "call"
	RETURN   = "return"
	KEY      = "key"
	META     = "meta"
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
