package cp

type ID int

type Id interface {
	Adr(...int) ID
}

var adr int = 0

func SomeAdr() ID {
	adr--
	return ID(adr)
}
