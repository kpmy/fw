package main

import (
	"os"
	"rt"
	"xev"
)

func main() {
	path, _ := os.Getwd()
	ret := xev.Load(path, "PrivDemo1.oxf")
	if ret != nil {
		p := rt.NewProcessor()
		err := p.ConnectTo(ret)
		if err != nil {
			panic("not connected")
		}
		for {
			res, _ := p.Do()
			if res != rt.OK {
				break
			}
		}
	} else {
		panic("no module")
	}
}
