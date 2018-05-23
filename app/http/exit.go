package http

import (
	"net/http"
	"strconv"

	"github.com/hailongz/kk-room/app"
)

func Exit(container app.IContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		data := GetInputData(r)

		id := data["id"]

		if id == "" {
			SetOutputData(map[string]interface{}{"errno": 100, "errmsg": "未找到应用ID"}, w)
			return
		}

		iid, _ := strconv.ParseInt(id, 10, 64)

		container.Exit(iid)

		SetOutputData(map[string]interface{}{"id": iid, "errno": 200}, w)

	}
}
