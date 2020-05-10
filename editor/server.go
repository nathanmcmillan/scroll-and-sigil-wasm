package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	contentType = "Content-type"
	textPlain   = "text/plain"
	dir         = "./public"
	api         = "/api"
	home        = dir + "/editor.html"
)

var extensions = map[string]string{
	".html": "text/html",
	".js":   "text/javascript",
	".css":  "text/css",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".wav":  "audio/wav",
	".mp3":  "audio/mpeg",
	".json": "application/json",
	".ttf":  "application/font-ttf",
	".wasm": "application/wasm",
}

func main() {
	num := len(os.Args)

	port := "3000"
	if num > 1 {
		port = os.Args[1]
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	httpserver := &http.Server{Addr: ":" + port, Handler: http.HandlerFunc(editor)}
	fmt.Println("listening on port " + port)

	go func() {
		err := httpserver.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-stop
	fmt.Println("signal interrupt")
	httpserver.Shutdown(context.Background())
	fmt.Println()
}

func editor(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, r.Method, r.URL.Path)

	if r.URL.Path == "/map" && r.Method == "POST" {
		w.Header().Set(contentType, textPlain)
		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		content := string(raw)
		index := strings.Index(content, ":")
		name := content[:index]
		data := content[index+1:]

		file, err := os.Create(filepath.Join("maps", name+".map"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString(data)
		if err != nil {
			panic(err)
		}

		return
	}

	var path string
	if r.URL.Path == "/" {
		path = home
	} else if strings.HasSuffix(r.URL.Path, ".map") {
		path = "maps" + r.URL.Path
	} else {
		path = dir + r.URL.Path
	}

	file, err := os.Open(path)
	if err != nil {
		path = home
		file, err = os.Open(path)
		if err != nil {
			return
		}
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	typ, has := extensions[filepath.Ext(path)]
	if !has {
		typ = textPlain
	}

	w.Header().Set(contentType, typ)
	w.Write(contents)
}
