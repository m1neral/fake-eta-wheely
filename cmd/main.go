package main

import (
	"log"
	"net/http"
	"os"

	eta "github.com/m1neral/fake-eta-wheely"
)

const (
	path = "/eta"
	port = ":8080"
)

func main() {
	apiRequestService := eta.NewApiRequestService()
	etaService := eta.NewEtaService(apiRequestService)

	logger := log.New(os.Stdout, "", 0)
	handler := eta.MakeHandler(etaService, logger)

	http.HandleFunc(path, handler)
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
