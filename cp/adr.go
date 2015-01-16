package cp

type Id interface {
	Adr(...int) int
}

var adr int = 0

func SomeAdr() int {
	adr--
	return adr
}
