package handler

import (
	"encoding/json"
	"net/http"
)

type HelloWorld struct {
	Message string `json:"message"`
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	msg := HelloWorld{Message: "Hello World"}
	jsonResp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}
