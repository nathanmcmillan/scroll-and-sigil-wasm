package main

import (
	"syscall/js"
)

type domJS struct {
	doc     js.Value
	window  js.Value
	body    js.Value
	console js.Value
}

var (
	dom *domJS
)

func domInit() *domJS {
	dom := &domJS{}
	dom.doc = js.Global().Get("document")
	dom.window = js.Global().Get("window")
	dom.body = dom.doc.Get("body")
	dom.console = js.Global().Get("console")
	return dom
}

func main() {
	done := make(chan struct{})
	dom = domInit()
	app := appInit()
	app.run()
	<-done
}

func console(c ...interface{}) {
	dom.console.Call("log", c...)
}
