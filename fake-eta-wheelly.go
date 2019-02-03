package main

import (
	"log"
	"net/http"
	"os"
)

const (
	path = "/eta"
	port = "8080"
)

func main() {
	apiRequestService := NewApiRequestService()
	etaService := NewEtaService(apiRequestService)

	logger := log.New(os.Stdout, "", 0)
	handler := MakeHandler(etaService, logger)

	http.HandleFunc(path, handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
