package main

import (
	"syscall/js"
)

var (
	inputKeys = map[string]bool{}
)

// InputIsKeyDown func
func InputIsKeyDown(key string) bool {
	return inputKeys[key]
}

// InputIsKeyPress func
func InputIsKeyPress(key string) bool {
	temp := inputKeys[key]
	inputKeys[key] = false
	return temp
}

// InputSetKeyUp func
func InputSetKeyUp() js.Func {
	return js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		inputKeys[args[0].Get("key").String()] = false
		return nil
	})
}

// InputSetKeyDown func
func InputSetKeyDown() js.Func {
	return js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		inputKeys[args[0].Get("key").String()] = true
		return nil
	})
}
