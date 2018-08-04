package ws

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

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
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("未找到房间ID"))
			return
		}

		iid, _ := strconv.ParseInt(roomId, 10, 64)

		R := server.RoomGet(iid)

		if R == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("未找到房间ID"))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("[" + r.RemoteAddr + "] [OPEN]")

		id := server.AutoId()

		log.Println("["+r.RemoteAddr+"] [ID]", id)

		ch := room.NewWSChannel(id, conn, 20480)

		R.AddChannel(ch)

		defer ch.Close()

		for {

			mType, data, err := conn.ReadMessage()

			if err != nil {
				log.Printf("[%s] [%d] [ERROR] %s\n", r.RemoteAddr, id, err.Error())
				break
			}

			if mType != websocket.BinaryMessage {
				log.Printf("[%s] [%d] [ERROR] %s\n", r.RemoteAddr, id, "Message Type Not Is Binary")
				break
			}

			message := kk.Message{}

			err = proto.Unmarshal(data, &message)

			if err != nil {
				log.Printf("[%s] [%d] [ERROR] %s\n", r.RemoteAddr, id, err.Error())
				break
			}

			if message.Type == kk.Message_PING {

				message.RoomId = iid
				message.From = id
				message.Dtime = (time.Now().UnixNano() / 1000000)
				message.Type = kk.Message_PONG
				err = ch.SendMessage(&message)

				if err != nil {
					log.Printf("[%s] [%d] [ERROR] %s\n", r.RemoteAddr, id, err.Error())
					break
				}

			} else {

				message.RoomId = iid
				message.Dtime = (time.Now().UnixNano() / 1000000)
				R.Send(&message)

			}
		}

		R.RemoveChannel(ch)

		log.Printf("[%s] [%d] [CLOSE]\n", r.RemoteAddr, id)

	}
}
