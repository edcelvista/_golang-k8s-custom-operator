package routes

import (
	controller "_gorestapi-k8s/controller"

	"github.com/gorilla/mux"
)

func (m *Router) WebhookRoutes() *mux.Router {
	m.R.HandleFunc("/webhook/{name}", controller.WebhookHandlerGET).Methods("GET")
	m.R.HandleFunc("/webhook/validating/pod", controller.WebhookValidatingHandlerPOSTPod).Methods("POST")
	m.R.HandleFunc("/webhook/mutating/pod", controller.WebhookMutatingHandlerPOSTPod).Methods("POST")
	m.R.HandleFunc("/webhook/validating/tenant", controller.WebhookValidatingHandlerPOSTTenant).Methods("POST")
	return m.R
}
