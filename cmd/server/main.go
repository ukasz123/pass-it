package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"pass-it/internal/cache"
	"pass-it/internal/server/handlers"
	"pass-it/internal/server/models"
)

func main() {
	var cache = cache.NewCache[models.DefaultStoredData]()
	var payloadMessageChannel = make(chan models.PayloadMessage[string, string])

	var storeRequest = handlers.NewStoreRequestHandler(cache)
	var fetchRequest = handlers.NewFetchRequestHandler(cache, payloadMessageChannel)
	var confirmRequest = handlers.NewConfirmRequestHandler(cache, payloadMessageChannel)
	var indexRequest = handlers.NewIndexRequestHandler()

	r := mux.NewRouter()
	r.Handle("/", indexRequest).Methods(http.MethodGet)
	r.Handle("/store/{id}", storeRequest).Methods(http.MethodPut)
	r.Handle("/store/{id}", confirmRequest).Methods(http.MethodPost)
	r.Handle("/fetch", fetchRequest).Methods(http.MethodGet)

	http.ListenAndServe(":8080", r)
}
