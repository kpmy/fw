package context

type Domain interface {
	ConnectTo(c ContextAware)
}

type ContextAware interface {
	Init(d Domain)
	Domain() Domain
}
