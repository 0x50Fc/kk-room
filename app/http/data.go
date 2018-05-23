package http

import (
	"encoding/json"
	"net/http"
)

func GetInputData(r *http.Request) map[string]string {

	v := map[string]string{}

	{
		vs := r.URL.Query()

		for key, values := range vs {
			v[key] = values[0]
		}

	}

	if r.Method == "POST" {

		r.ParseForm()

		for key, values := range r.Form {
			v[key] = values[0]
		}
	}

	return v
}

func SetOutputData(data interface{}, w http.ResponseWriter) {
	b, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
