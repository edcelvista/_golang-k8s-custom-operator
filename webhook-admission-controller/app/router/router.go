package router

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	Controller "_gorestapi-k8s/controller"
)

type Router struct {
	r *mux.Router
}

func (m *Router) PingRoutes() *mux.Router {
	m.r.HandleFunc("/ping/{name}", Controller.PingHandlerGET).Methods("GET")
	m.r.HandleFunc("/ping", Controller.PingHandlerPOST).Methods("POST")
	return m.r
}

func (m *Router) WebhookRoutes() *mux.Router {
	m.r.HandleFunc("/webhook/{name}", Controller.WebhookHandlerGET).Methods("GET")
	m.r.HandleFunc("/webhook/validating/{name}", Controller.WebhookValidatingHandlerPOST).Methods("POST")
	m.r.HandleFunc("/webhook/mutating/{name}", Controller.WebhookMutatingHandlerPOST).Methods("POST")
	return m.r
}

func Run() {
	// Load .env file
	err := godotenv.Load()
	port := ":8443"

	// TLS config
	cert := "./tls.crt"
	key := "./tls.key"

	if err != nil {
		log.Println("⚠️ Error loading .env file")
		// Get environment variables
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
		if os.Getenv("SSL_CERT") != "" {
			cert = os.Getenv("SSL_CERT")
		}
		if os.Getenv("SSL_KEY") != "" {
			key = os.Getenv("SSL_KEY")
		}
		log.Printf("⚠️ Defaulting to Port %v", port)
	} else {
		port, _ = os.LookupEnv("PORT")
		cert, _ = os.LookupEnv("SSL_CERT")
		key, _ = os.LookupEnv("SSL_KEY")
		log.Println("💡 Found .env")
	}

	muxRouter := mux.NewRouter()
	router := Router{
		r: muxRouter,
	}

	router.PingRoutes()
	router.WebhookRoutes()

	http.Handle("/", muxRouter)

	// TLS config
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// TLS config
	server := &http.Server{
		Addr:      port,
		Handler:   muxRouter,
		TLSConfig: cfg,
	}

	log.Printf("💡 ⚡️ Mux API Running @ %s \n", port)
	// err = http.ListenAndServe(port, muxRouter)

	// TLS config
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("‼️ Failed to start router %s with %v %v", err, cert, key)
	}
}
