package main

import (
	"flag"
	"fmt"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/module"
	"fw/rt2/nodeframe"
	_ "fw/rt2/rules"
	"fw/rt2/scope"
	_ "fw/rt2/scope/std"
	"time"
	"ypk/assert"
)

var name string

func init() {
	flag.StringVar(&name, "i", "", "-i name.ext")
}
func close() {
	fmt.Println("closed")
}

func main() {
	flag.Parse()
	if name == "" {
		name = "XevDemo15"
	}
	global := new(stdDomain)
	modList := module.New()
	global.ConnectTo(context.MOD, modList)
	t0 := time.Now()
	ret, err := modList.Load(name)
	t1 := time.Now()
	fmt.Println("load", t1.Sub(t0))
	assert.For(ret != nil, 40)
	assert.For(err == nil, 41)
	defer close()
	{
		domain := new(stdDomain)
		global.ConnectTo(name, domain)
		root := frame.NewRoot()
		domain.ConnectTo(context.STACK, root)
		domain.ConnectTo(context.SCOPE, scope.New())
		var fu nodeframe.FrameUtils
		root.PushFor(fu.New(ret.Enter), nil)
		i := 0
		t0 := time.Now()
		for x := frame.NOW; x == frame.NOW; x = root.Do() {
			//fmt.Println(x)
			i++
		}
		t1 := time.Now()
		fmt.Println("total steps", i)
		fmt.Println("spent", t1.Sub(t0))
	}
}
