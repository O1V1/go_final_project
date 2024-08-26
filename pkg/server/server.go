package server

import (
	"log"
	"net/http"

	"github.com/O1V1/go_final_project/pkg/config"
)

func Start() {
	log.Printf("Starting server on port %s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
