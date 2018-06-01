package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hailongz/kk-room/app"
)

func Open(container app.IContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := GetInputData(r)

		path := data["path"]

		if path == "" {
			SetOutputData(map[string]interface{}{"errno": 101, "errmsg": "未找到应用路径"}, w)
			return
		}

		ch := make(chan bool)

		id := data["id"]
		iid, _ := strconv.ParseInt(id, 10, 64)
		expires, _ := strconv.ParseInt(data["expires"], 10, 64)

		delete(data, "expires")
		delete(data, "path")

		iid = container.Open(iid, path, data, time.Duration(expires)*time.Second, ch)

		close(ch)

		SetOutputData(map[string]interface{}{"id": iid}, w)

	}
}
