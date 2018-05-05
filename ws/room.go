package ws

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hailongz/kk-room/proto/golang/kk"
	"github.com/hailongz/kk-room/room"
)

func Room(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		roomId := query.Get("id")

		if roomId == "" {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			return
		}

		iid, _ := strconv.ParseInt(roomId, 10, 64)

		R := server.RoomGet(iid)

		if R == nil {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("Connected: ", r.RemoteAddr)

		id := server.AutoId()

		ch := room.NewWSChannel(id, c, 204800)

		R.AddChannel(ch)

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
					message.RoomId = iid
					message.Dtime = (time.Now().UnixNano() / 1000000)
					message.Type = kk.Message_PONG
					ch.Send(message)
				} else {

					message.RoomId = iid
					message.Dtime = (time.Now().UnixNano() / 1000000)
					R.Send(message)

				}

			case <-ch.CloseChannel():
				log.Printf("Disconnected: %s %s\n", r.RemoteAddr, ch.GetError())
				run = false
				break
			}
		}

		R.RemoveChannel(ch)

	}
}
