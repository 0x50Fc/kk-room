package admin

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hailongz/kk-room/room"
)

func RoomCreate(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := map[string]interface{}{}

		inputData := GetInputData(r)

		id := inputData["id"]

		iid, _ := strconv.ParseInt(id, 10, 64)

		expires, _ := strconv.ParseInt(inputData["expires"], 10, 64)

		room := server.RoomCreate(iid, time.Duration(expires)*time.Second)

		data["id"] = room.GetId()
		data["errno"] = 200
		data["errmsg"] = "成功"

		SetOutputData(data, w)

	}
}

func RoomExit(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)
		data := map[string]interface{}{}

		inputData := GetInputData(r)

		id, ok := inputData["id"]

		for {

			if !ok || id == "" {
				data["errno"] = 400
				data["errmsg"] = "未找到房间ID"
				break
			}

			iid, err := strconv.ParseInt(id, 10, 64)

			if err != nil {
				data["errno"] = 401
				data["errmsg"] = "错误的房间ID"
				break
			}

			server.RoomRemove(iid)

			data["id"] = iid
			data["errno"] = 200
			data["errmsg"] = "成功"

			break
		}

		SetOutputData(data, w)
	}
}

func RoomGet(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)
		data := map[string]interface{}{}

		inputData := GetInputData(r)

		id, ok := inputData["id"]

		for {

			if !ok || id == "" {
				data["errno"] = 400
				data["errmsg"] = "未找到房间ID"
				break
			}

			iid, err := strconv.ParseInt(id, 10, 64)

			if err != nil {
				data["errno"] = 401
				data["errmsg"] = "错误的房间ID"
				break
			}

			room := server.RoomGet(iid)

			if room == nil {
				data["errno"] = 404
				data["errmsg"] = "未找到房间"
				break
			}

			data["id"] = room.GetId()
			data["errno"] = 200
			data["errmsg"] = "成功"

			break
		}

		SetOutputData(data, w)
	}
}

func RoomList(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := map[string]interface{}{}

		data["items"] = server.RoomList()
		data["errno"] = 200
		data["errmsg"] = "成功"

		SetOutputData(data, w)
	}
}
