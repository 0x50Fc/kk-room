package http

import (
	"log"
	"net/http"

	"github.com/hailongz/kk-room/app"
)

func List(container app.IContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		ch := make(chan bool)

		ids := container.List(ch)

		close(ch)

		SetOutputData(map[string]interface{}{"apps": ids, "errno": 200}, w)

	}
}
