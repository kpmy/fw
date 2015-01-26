package main

import (
	"flag"
	"fmt"
	mod "fw/cp/module"
	"fw/rt2/context"
	"fw/rt2/decision"
	rtmod "fw/rt2/module"
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
}

func main() {
	flag.Parse()
	if name == "" {
		name = "XevDemo19"
	}
	global := &stdDomain{god: true}
	global.global = global
	modList := rtmod.New()
	global.Attach(context.MOD, modList)
	heap = scope.New(context.HEAP)
	global.Attach(context.HEAP, heap)
	t0 := time.Now()
	var init []*mod.Module
	_, err := modList.Load(name, func(m *mod.Module) {
		init = append(init, m)
	})
	t1 := time.Now()
	fmt.Println("load", t1.Sub(t0))
	assert.For(err == nil, 40)
	defer close()
	decision.Run(global, init)
}
