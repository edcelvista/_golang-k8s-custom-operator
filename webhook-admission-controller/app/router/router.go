package router

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	Controller "_gorestapi-k8s/controller"
	Lib "_gorestapi-k8s/lib"
)

type Router struct {
	r *mux.Router
}

type CertsAndKeys struct {
	cert string
	key  string
}

func (m *CertsAndKeys) checkCerts() {
	// Load certificate
	certPath := m.cert
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatalf("‼️ Error reading certificate: %v %v", err, m.cert)
	}

	// Decode the PEM data
	block, _ := pem.Decode(certData)
	if block == nil {
		log.Fatalf("Failed to decode PEM block")
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("‼️ Error parsing certificate: %v %v", err, m.cert)
	}

	// Check if the certificate is expired
	currentTime := time.Now()
	if currentTime.After(cert.NotAfter) {
		log.Println("‼️ Certificate has expired.")
	} else {
		log.Println("💡 Certificate is valid.")
	}

	// Check the NotBefore field (if the certificate is not yet valid)
	if currentTime.Before(cert.NotBefore) {
		log.Println("‼️ Certificate is not yet valid.")
	} else {
		log.Println("💡 Certificate is within the valid period.")
	}

	// Optionally, you could also validate other parts, like the issuer and subject
	log.Println("🔑 Issuer:", cert.Issuer)
	log.Println("🔑 Subject:", cert.Subject)
	log.Println("🔑 Not After:", cert.NotAfter)
	log.Println("🔑 Not Before:", cert.NotBefore)
}

func (m *Router) PingRoutes() *mux.Router {
	m.r.HandleFunc("/ping/{name}", Controller.PingHandlerGET).Methods("GET")
	m.r.HandleFunc("/ping", Controller.PingHandlerPOST).Methods("POST")
	m.r.HandleFunc("/echo", Controller.EchoHandlerPOST).Methods("POST")
	return m.r
}

func (m *Router) WebhookRoutes() *mux.Router {
	m.r.HandleFunc("/webhook/{name}", Controller.WebhookHandlerGET).Methods("GET")
	m.r.HandleFunc("/webhook/validating/pod", Controller.WebhookValidatingHandlerPOSTPod).Methods("POST")
	m.r.HandleFunc("/webhook/mutating/pod", Controller.WebhookMutatingHandlerPOSTPod).Methods("POST")
	m.r.HandleFunc("/webhook/validating/tenant", Controller.WebhookValidatingHandlerPOSTTenant).Methods("POST")
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

	tlsCertsAndKey := CertsAndKeys{
		cert: cert,
		key:  key,
	}

	tlsCertsAndKey.checkCerts()
	Lib.Init()

	log.Printf("💡 ⚡️ Mux API Running 📦 %s with 🔑 %v %v \n", port, cert, key)
	// err = http.ListenAndServe(port, muxRouter)

	// TLS config
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("‼️ Failed to start router %s with %v %v", err, cert, key)
	}

}
