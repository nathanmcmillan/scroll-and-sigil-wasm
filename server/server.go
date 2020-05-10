package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	contentType = "Content-type"
	textPlain   = "text/plain"
	dir         = "./public"
	home        = dir + "/game.html"
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

var (
	secure = false
)

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	num := len(os.Args)

	port := "3000"
	if num > 1 {
		port = os.Args[1]
	}

	level := "maps/test.map"
	if num > 2 {
		level = "maps/" + os.Args[2] + ".map"
	}

	serveFunction := game(level)
	httpserver := &http.Server{
		Addr:         ":" + port,
		Handler:      http.HandlerFunc(serveFunction),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if num > 3 && os.Args[3] == "-secure" {
		secure = true
		cert := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("scrollandsigil.eastus.cloudapp.azure.com"),
			Cache:      autocert.DirCache("certs"),
		}
		httpserver.TLSConfig = &tls.Config{GetCertificate: cert.GetCertificate}

		fmt.Println("listening on port " + port + " (https)")
		go func() {
			err := http.ListenAndServe(":http", cert.HTTPHandler(nil))
			if err != nil {
				fmt.Println(err)
			}
		}()
		go func() {
			err := httpserver.ListenAndServeTLS("", "")
			if err != nil {
				fmt.Println(err)
			}
		}()

	} else {
		fmt.Println("listening on port " + port)
		go func() {
			err := httpserver.ListenAndServe()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	<-stop
	fmt.Println("signal interrupt")
	httpserver.Shutdown(context.Background())
	fmt.Println()
}
