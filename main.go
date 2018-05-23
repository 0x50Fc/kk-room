package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/hailongz/kk-room/admin"
	"github.com/hailongz/kk-room/app"
	H "github.com/hailongz/kk-room/app/http"
	"github.com/hailongz/kk-room/room"
	"github.com/hailongz/kk-room/ws"
)

var broadcast_addr = flag.String("broadcast", "", "广播服务地址 如 :8080")
var admin_addr = flag.String("admin", "", "广播管理地址 如 :8081")
var app_addr = flag.String("app", "", "应用管理地址 如 :8082")
var server room.IServer = nil
var container app.IContainer = nil

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if *broadcast_addr != "" {

		if server == nil {
			server = room.NewServer(204800)
		}

		go func() {

			log.Printf("[BROADCAST] %s\n", *broadcast_addr)

			mux := http.NewServeMux()

			mux.HandleFunc("/room", ws.Room(server))

			s := &http.Server{
				Addr:           *broadcast_addr,
				Handler:        mux,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			log.Fatal(s.ListenAndServe())
		}()

		if *admin_addr != "" {

			go func() {

				log.Printf("[ADMIN] %s\n", *admin_addr)

				mux := http.NewServeMux()

				mux.HandleFunc("/room/create", admin.RoomCreate(server))
				mux.HandleFunc("/room/get", admin.RoomGet(server))
				mux.HandleFunc("/room/exit", admin.RoomExit(server))
				mux.HandleFunc("/room/list", admin.RoomList(server))
				mux.HandleFunc("/runtime/state", admin.RuntimeState(server))
				mux.HandleFunc("/runtime/conns", admin.RuntimeConns(server))

				s := &http.Server{
					Addr:           *admin_addr,
					Handler:        mux,
					ReadTimeout:    10 * time.Second,
					WriteTimeout:   10 * time.Second,
					MaxHeaderBytes: 1 << 20,
				}

				log.Fatal(s.ListenAndServe())
			}()

		}
	}

	if *app_addr != "" {

		container = app.NewContainer()

		go func() {

			log.Printf("[APP] %s\n", *app_addr)

			mux := http.NewServeMux()

			mux.HandleFunc("/app/open", H.Open(container))
			mux.HandleFunc("/app/get", H.Get(container))
			mux.HandleFunc("/app/exit", H.Exit(container))
			mux.HandleFunc("/app/command", H.Command(container))
			mux.HandleFunc("/app/list", H.List(container))

			s := &http.Server{
				Addr:           *app_addr,
				Handler:        mux,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			log.Fatal(s.ListenAndServe())
		}()

	}

	if *broadcast_addr != "" || *app_addr != "" {
		ch := make(chan bool)
		<-ch
		close(ch)
	} else {
		flag.PrintDefaults()
	}

}
