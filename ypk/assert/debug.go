package assert

func For(cond bool, code int) {
	if !cond {
		panic(code)
	}
}
