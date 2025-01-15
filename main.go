package main

import (
	"log"
	"net/http"

	"github.com/MunifTanjim/stremthru/internal/config"
	"github.com/MunifTanjim/stremthru/internal/db"
	"github.com/MunifTanjim/stremthru/internal/endpoint"
)

func main() {
	config.PrintConfig()

	mux := http.NewServeMux()

	endpoint.AddRootEndpoint(mux)
	endpoint.AddHealthEndpoints(mux)
	endpoint.AddProxyEndpoints(mux)
	endpoint.AddStoreEndpoints(mux)
	endpoint.AddStremioEndpoints(mux)

	database := db.Open()
	defer db.Close()
	db.Ping()
	RunSchemaMigration(database.URI)

	addr := ":" + config.Port
	server := &http.Server{Addr: addr, Handler: mux}

	server.SetKeepAlivesEnabled(false)

	log.Println("stremthru listening on " + addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start stremthru: %v", err)
	}
}
