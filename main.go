package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/hailongz/kk-room/admin"
	"github.com/hailongz/kk-room/room"
	"github.com/hailongz/kk-room/ws"
)

var addr = flag.String("addr", ":8080", "ws port")
var arg_admin = flag.String("admin", ":8081", "admin port")
var server room.IServer = nil

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	server = room.NewServer(2048)

	if *arg_admin != "" {

		go func() {

			log.Printf("Room Admin %s\n", *arg_admin)

			mux := http.NewServeMux()

			mux.HandleFunc("/room/create", admin.RoomCreate(server))
			mux.HandleFunc("/room/get", admin.RoomGet(server))
			mux.HandleFunc("/room/exit", admin.RoomExit(server))
			mux.HandleFunc("/runtime/state", admin.RuntimeState(server))
			mux.HandleFunc("/runtime/conns", admin.RuntimeConns(server))

			s := &http.Server{
				Addr:           *arg_admin,
				Handler:        mux,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			log.Fatal(s.ListenAndServe())
		}()

	}

	http.HandleFunc("/room", ws.Room(server))
	log.Printf("Room Server %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
