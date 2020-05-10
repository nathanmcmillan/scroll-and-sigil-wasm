package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Server struct
type Server struct {
	world  *World
	people []*Person
	mux    *sync.Mutex
}

var (
	server *Server
)

func game(level string) func(w http.ResponseWriter, r *http.Request) {
	server = &Server{}
	server.mux = &sync.Mutex{}
	server.people = make([]*Person, 0)

	server.world = NewWorld()
	file, err := os.Open(level)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	server.world.Load(contents)

	go func() {
		const rate = time.Duration(WorldTickRate) * time.Millisecond
		ticker := time.Now().Add(rate)
		for true {
			server.mux.Lock()
			num := len(server.people)
			if num > 0 {
				server.world.Update()
				server.world.BuildSnapshots(server.people)
			}
			server.mux.Unlock()
			time.Sleep(ticker.Sub(time.Now()))
			ticker = ticker.Add(rate)
			server.mux.Lock()
			num = len(server.people)
			for i := 0; i < num; i++ {
				person := server.people[i]
				if person.snap != nil {
					go person.WriteBinaryToClient(person.snap)
					person.snap = nil
				}
			}
			server.mux.Unlock()
		}
	}()

	return serve
}

func serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, r.Method, r.URL.Path)

	var path string
	if r.URL.Path == "/" {
		path = home
	} else if r.URL.Path == "/websocket" {
		server.connectSocket(w, r)
		return
	} else {
		path = dir + r.URL.Path
	}

	file, err := os.Open(path)
	defer file.Close()
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

func (me *Server) connectSocket(writer http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	safe := false
	if secure {
		if origin == "https://"+request.Host {
			safe = true
		}
	} else if origin == "http://"+request.Host {
		safe = true
	}
	if !safe {
		http.Error(writer, "origin not allowed", 403)
		return
	}
	// TODO client's need to be given a key to avoid DOS attacks / wasting resources on non players
	// TODO http server read list of acceptable path to files preventing attacks
	// TODO editor should be publicly accesible, facilitates community and longevity
	upgrader := websocket.Upgrader{}
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, "could not open websocket", 400)
		return
	}
	me.mux.Lock()
	person := NewPerson(connection, server.world)
	me.people = append(me.people, person)
	data := me.world.Save(person)
	me.world.broadcastCount++
	me.world.broadcast.PutUint8(BroadcastNew)
	person.Character.Save(me.world.broadcast)
	me.mux.Unlock()
	person.WriteBinaryToClient(data)
	go me.PersonConnectionLoop(person)
}

// PersonConnectionLoop func
func (me *Server) PersonConnectionLoop(person *Person) {
	for {
		_, data, err := person.Connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			person.Connection.Close()
			break
		}
		me.mux.Lock()
		person.Input(data)
		me.mux.Unlock()
	}

	char := person.Character
	char.Health = 0
	char.World.removeThing(char.thing)
	char.removeFromBlocks()

	me.RemovePerson(person)
}

// RemovePerson func
func (me *Server) RemovePerson(person *Person) {
	me.mux.Lock()
	defer me.mux.Unlock()

	me.world.broadcastCount++
	me.world.broadcast.PutUint8(BroadcastDelete)
	me.world.broadcast.PutUint16(person.Character.NID)

	num := len(me.people)
	for i := 0; i < num; i++ {
		if me.people[i] == person {
			me.people[i] = me.people[num-1]
			me.people[num-1] = nil
			me.people = me.people[:num-1]
			return
		}
	}
}
