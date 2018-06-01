package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/hailongz/kk-room/app"
)

func Get(container app.IContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := GetInputData(r)

		id := data["id"]

		if id == "" {
			SetOutputData(map[string]interface{}{"errno": 100, "errmsg": "未找到应用ID"}, w)
			return
		}

		iid, _ := strconv.ParseInt(id, 10, 64)

		ch := make(chan bool)

		iid = container.Get(iid, ch)

		close(ch)

		if iid == 0 {
			SetOutputData(map[string]interface{}{"errno": 102, "errmsg": "未找到应用"}, w)
		} else {
			SetOutputData(map[string]interface{}{"id": iid, "errno": 200}, w)
		}

	}
}
