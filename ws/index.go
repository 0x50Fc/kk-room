package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hailongz/kk-room/proto/golang/kk"
	"github.com/hailongz/kk-room/room"
)

func Index(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		log.Println("Connected: ", r.RemoteAddr)

		var rom room.IRoom = nil

		id := server.AutoId()

		ch := room.NewWSChannel(id, c, 204800)

		defer ch.Close()

		run := true

		for run {
			select {
			case message := <-ch.ReadChannel():

				if message == nil {
					log.Printf("Disconnected: %s %s\n", r.RemoteAddr, ch.GetError())
					run = false
					break
				}

				if message.Type == kk.Message_PING {
					message.Dtime = (time.Now().UnixNano() / 1000000)
					message.Type = kk.Message_PONG
					ch.Send(message)
				} else if message.Type == kk.Message_FRAME {

					if rom == nil {
						roomId := message.RoomId
						rom = server.RoomGet(roomId)
						if rom != nil {
							rom.AddChannel(ch)
						}
					}

					if rom != nil {
						message.Dtime = (time.Now().UnixNano() / 1000000)
						rom.Send(message)
					}

				}

			case <-ch.CloseChannel():
				log.Printf("Disconnected: %s %s\n", r.RemoteAddr, ch.GetError())
				run = false
				break
			}
		}

		if rom != nil {
			rom.RemoveChannel(ch)
		}

	}
}
