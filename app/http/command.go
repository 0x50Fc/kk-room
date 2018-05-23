package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hailongz/kk-room/app"
)

func Command(container app.IContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		data := GetInputData(r)

		id := data["id"]

		if id == "" {
			SetOutputData(map[string]interface{}{"errno": 100, "errmsg": "未找到应用ID"}, w)
			return
		}

		iid, _ := strconv.ParseInt(id, 10, 64)

		b, _ := json.Marshal(data)

		container.RunCommand(iid, b)

		if iid == 0 {
			SetOutputData(map[string]interface{}{"errno": 102, "errmsg": "未找到应用"}, w)
		} else {
			SetOutputData(map[string]interface{}{"id": iid, "errno": 200}, w)
		}

	}
}
