package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"universal-response/home"
	"universal-response/server"
)

var (
	CertFile       = os.Getenv("UNIVERSAL_RESPONSE_CERT_FILE")
	KeyFile        = os.Getenv("UNIVERSAL_RESPONSE_KEY_FILE")
	ServiceAddr    = os.Getenv("UNIVERSAL_RESPONSE_SERVICE_ADDR")
	ServiceSslAddr = os.Getenv("UNIVERSAL_RESPONSE_SERVICE_SSL_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	mux := http.NewServeMux()
	home.SetupRoutes(mux)

	errors := make(chan error)
	var ssl map[string]string
	if ServiceSslAddr != "" {
		ssl = make(map[string]string)
		ssl["key"] = KeyFile
		ssl["cert"] = CertFile
		logger.Println(fmt.Sprintf("Server starting (%s, %s)", ServiceAddr, ServiceSslAddr))
	} else {
		logger.Println(fmt.Sprintf("Server starting (%s)", ServiceAddr))
	}

	errors = server.StartServer(mux, ServiceAddr, ServiceSslAddr, ssl)

	select {
	case err := <-errors:
		logger.Fatalf("server failed to start: %v", err)
	}
}
