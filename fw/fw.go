package main

import (
	"fmt"
	"os"
	"rt2/context"
	"rt2/frame"
	"rt2/nodeframe"
	_ "rt2/rules"
	"xev"
	"ypk/assert"
)

type stdDomain struct {
}

func (d *stdDomain) ConnectTo(x context.ContextAware) {
	assert.For(x != nil, 20)
	x.Init(d)
}

func main() {
	path, _ := os.Getwd()
	ret := xev.Load(path, "PrivDemo1.oxf")
	assert.For(ret != nil, 20)
	domain := new(stdDomain)
	root := frame.NewRoot()
	domain.ConnectTo(root)
	var fu nodeframe.FrameUtils
	root.Push(fu.New(ret.Enter))
	i := 0
	for x := frame.DO; x == frame.DO; x = root.Do() {
		fmt.Println(x)
		i++
	}
	fmt.Println("total steps", i)
}
