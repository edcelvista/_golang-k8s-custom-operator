package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	Model "_gorestapi-k8s/model"
)

func PingHandlerGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	res := Model.Pong{Pong: "pong", Message: reqParams}
	json.NewEncoder(w).Encode(res)
}

func PingHandlerPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var parsedUserRequest Model.Ping
	json.NewDecoder(r.Body).Decode(&parsedUserRequest)

	res := Model.Pong{Pong: "pong", Message: parsedUserRequest.Message}
	json.NewEncoder(w).Encode(res)
}
