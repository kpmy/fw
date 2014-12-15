package main

import (
	"fmt"
	"os"
	"rt2/frame"
	"rt2/nodeframe"
	_ "rt2/rules"
	"xev"
	"ypk/assert"
)

func main() {
	path, _ := os.Getwd()
	ret := xev.Load(path, "PrivDemo1.oxf")
	assert.For(ret != nil, 20)
	root := new(frame.RootFrame).Init()
	var fu nodeframe.FrameUtils
	root.Push(fu.New(ret.Enter))
	i := 0
	for x := frame.DO; x == frame.DO; x = root.Do() {
		fmt.Println(x)
		i++
	}
	fmt.Println("total steps", i)
}
