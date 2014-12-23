package main

import (
	"fmt"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/module"
	"fw/rt2/nodeframe"
	_ "fw/rt2/rules"
	"fw/rt2/scope"
	"ypk/assert"
)

func main() {
	global := new(stdDomain)
	modList := module.New()
	global.ConnectTo(context.MOD, modList)
	ret, err := modList.Load("XevDemo5")
	assert.For(ret != nil, 40)
	assert.For(err == nil, 41)
	{
		domain := new(stdDomain)
		global.ConnectTo("XevDemo5", domain)
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
