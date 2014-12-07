package main

import (
	"fmt"
	"os"
	"rt"
	"xev"
)

func main() {
	path, _ := os.Getwd()
	ret := xev.Load(path, "PrivDemo1.oxf")
	if ret != nil {
		p := rt.NewProcessor()
		p.ConnectTo(ret)
		for {
			res, _ := p.Do()
			if res != rt.OK {
				break
			}
		}
		fmt.Println("")
	}
}
