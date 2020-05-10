package network

import (
	"sync"
	"syscall/js"
)

func fetch(promise js.Value) string {
	var text string
	var join sync.WaitGroup
	join.Add(1)
	success := js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		defer join.Done()

		join.Add(1)
		promise2 := args[0].Call("text")
		var success2 js.Func
		var fail2 js.Func
		success2 = js.FuncOf(func(self js.Value, args []js.Value) interface{} {
			defer success2.Release()
			defer fail2.Release()
			defer join.Done()
			text = args[0].String()
			return nil
		})
		fail2 = js.FuncOf(func(self js.Value, args []js.Value) interface{} {
			defer success2.Release()
			defer fail2.Release()
			defer join.Done()
			panic(args[0].String())
			return nil
		})
		promise2.Call("then", success2)
		promise2.Call("catch", fail2)

		return nil
	})
	fail := js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		defer join.Done()
		panic(args[0].String())
		return nil
	})
	promise.Call("then", success)
	promise.Call("catch", fail)
	join.Wait()
	success.Release()
	fail.Release()
	return text
}

// Get func
func Get(url string) string {
	location := js.Global().Get("location")
	origin := location.Get("origin").String()
	url = origin + "/" + url
	promise := js.Global().Call("fetch", url)
	return fetch(promise)
}

// Send func
func Send(url, data string) string {
	options := js.Global().Get("Object").New()
	options.Set("method", "POST")
	options.Set("body", data)
	location := js.Global().Get("location")
	origin := location.Get("origin").String()
	url = origin + "/" + url
	promise := js.Global().Call("fetch", url, options)
	return fetch(promise)
}

// Socket func
func Socket(url string) js.Value {
	location := js.Global().Get("location")
	host := location.Get("host").String()
	protocol := location.Get("protocol").String()
	url = host + "/" + url
	if protocol == "https:" {
		url = "wss://" + url
	} else {
		url = "ws://" + url
	}
	socket := js.Global().Get("WebSocket").New(url)
	err := js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		panic(args[0].String())
		return nil
	})
	socket.Set("onerror", err)
	return socket
}
