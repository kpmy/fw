package main

import (
	"flag"
	"fmt"
	"fw/cp"
	cpm "fw/cp/module"
	"fw/rt2/context"
	"fw/rt2/decision"
	rtm "fw/rt2/module"
	_ "fw/rt2/rules"
	"fw/rt2/scope"
	_ "fw/rt2/scope/modern"
	"fw/utils"
	"time"
	"ypk/assert"
)

var name string
var heap scope.Manager

func init() {
	flag.StringVar(&name, "i", "", "-i name.ext")
}

func close() {
	utils.PrintFrame("____")
	utils.PrintFrame(heap)
	utils.PrintFrame("^^^^")
	fmt.Println("closed")
	//fmt.Println(heap)
}

func main() {
	flag.Parse()
	if name == "" {
		name = "Start"
		utils.Debug(false)
	}
	global := &stdDomain{god: true}
	global.global = global
	modList := rtm.New()
	global.Attach(context.MOD, modList)
	global.Attach(context.DIGEST, context.Data(cp.Init()))
	heap = scope.New(context.HEAP)
	global.Attach(context.HEAP, heap)
	t0 := time.Now()
	var init []*cpm.Module
	_, err := modList.Load(name, func(m *cpm.Module) {
		init = append(init, m)
	})
	t1 := time.Now()
	fmt.Println("load", t1.Sub(t0))
	assert.For(err == nil, 40)
	defer close()
	decision.Run(global, init)
}
