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
	ret, err := modList.Load("XevDemo4")
	assert.For(ret != nil, 40)
	assert.For(err == nil, 41)
	{
		domain := new(stdDomain)
		global.ConnectTo("XevDemo4", domain)
		root := frame.NewRoot()
		domain.ConnectTo(context.STACK, root)
		domain.ConnectTo(context.SCOPE, scope.New())
		var fu nodeframe.FrameUtils
		root.PushFor(fu.New(ret.Enter), nil)
		i := 0
		for x := frame.NOW; x == frame.NOW; x = root.Do() {
			//fmt.Println(x)
			i++
		}
		fmt.Println("total steps", i)
	}
}
