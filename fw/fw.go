package main

import (
	"fmt"
	"rt2/context"
	"rt2/frame"
	"rt2/module"
	"rt2/nodeframe"
	_ "rt2/rules"
	"rt2/scope"
	"ypk/assert"
)

func main() {
	global := new(stdDomain)
	modList := module.New()
	global.ConnectTo(context.MOD, modList)
	ret, err := modList.Load("PrivDemo1")
	assert.For(err == nil, 20)
	{
		domain := new(stdDomain)
		global.ConnectTo("PrivDemo1", domain)
		root := frame.NewRoot()
		domain.ConnectTo(context.STACK, root)
		domain.ConnectTo(context.SCOPE, scope.New())
		var fu nodeframe.FrameUtils
		root.Push(fu.New(ret.Enter))
		i := 0
		for x := frame.DO; x == frame.DO; x = root.Do() {
			fmt.Println(x)
			i++
		}
		fmt.Println("total steps", i)
	}
}
