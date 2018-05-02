package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/hailongz/kk-room/room"
)

func RoomCreate(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := map[string]interface{}{}

		room := server.RoomCreate()

		data["roomId"] = room.GetId()
		data["errcode"] = 200
		data["errmsg"] = "成功创建房间"

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
				data["errcode"] = 400
				data["errmsg"] = "未找到房间ID"
				break
			}

			iid, err := strconv.ParseInt(id, 16, 64)

			if err != nil {
				data["errcode"] = 401
				data["errmsg"] = "错误的房间ID"
				break
			}

			server.RoomRemove(iid)

			data["errcode"] = 200
			data["errmsg"] = "成功结束房间"

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
				data["errcode"] = 400
				data["errmsg"] = "未找到房间ID"
				break
			}

			iid, err := strconv.ParseInt(id, 16, 64)

			if err != nil {
				data["errcode"] = 401
				data["errmsg"] = "错误的房间ID"
				break
			}

			room := server.RoomGet(iid)

			if room == nil {
				data["errcode"] = 404
				data["errmsg"] = "未找到房间"
				break
			}

			data["roomId"] = room.GetId()
			data["errcode"] = 200
			data["errmsg"] = "成功结束房间"

			break
		}

		SetOutputData(data, w)
	}
}
