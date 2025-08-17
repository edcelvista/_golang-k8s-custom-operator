package routes

import (
	controller "_gorestapi-k8s/controller"

	"github.com/gorilla/mux"
)

func (m *Router) PingRoutes() *mux.Router {
	m.R.HandleFunc("/healthz", controller.HealthHandlerGET).Methods("GET")
	m.R.HandleFunc("/ping/{name}", controller.PingHandlerGET).Methods("GET")
	m.R.HandleFunc("/ping", controller.PingHandlerPOST).Methods("POST")
	m.R.HandleFunc("/echo", controller.EchoHandlerPOST).Methods("POST")
	return m.R
}
