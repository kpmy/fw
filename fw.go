package main

import (
	"flag"
	"fw/cp"
	cpm "fw/cp/module"
	"fw/rt2/context"
	"fw/rt2/decision"
	rtm "fw/rt2/module"
	_ "fw/rt2/rules2"
	"fw/rt2/scope"
	_ "fw/rt2/scope/modern"
	"fw/utils"
	"log"
	"time"
	"ypk/assert"
)

var name string
var debug bool = false
var heap scope.Manager

func init() {
	flag.StringVar(&name, "i", "", "-i name.ext")
	flag.BoolVar(&debug, "d", false, "-d true/false")
}

func close() {
	utils.Debug(false)
	utils.PrintFrame("____")
	utils.PrintFrame(heap)
	utils.PrintFrame("^^^^")
	log.Println("closed")
}

func main() {
	flag.Parse()
	if name == "" {
		name = "Start3"
		utils.Debug(true)
	}
	utils.Debug(debug)
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
	log.Println("load", t1.Sub(t0))
	assert.For(err == nil, 40)
	defer close()
	decision.Run(global, init)
}
