package main

import (
	"net/http"
	"pass-it/cache"
	"pass-it/server/handlers"
	"pass-it/server/models"
	"github.com/gorilla/mux"
)

func main() {
	var cache = cache.NewCache[models.DefaultStoredData]()
	var payloadMessageChannel = make(chan models.PayloadMessage[string, string])

	var storeRequest = handlers.NewStoreRequestHandler(cache)
	var fetchRequest = handlers.NewFetchRequestHandler(cache, payloadMessageChannel)
	var confirmRequest = handlers.NewConfirmRequestHandler(cache, payloadMessageChannel)

	r := mux.NewRouter()
	r.Handle("/store/{id}",storeRequest).Methods(http.MethodPut)
	r.Handle("/store/{id}",confirmRequest).Methods(http.MethodPost)
	r.Handle("/fetch", fetchRequest).Methods(http.MethodGet)
	

	http.ListenAndServe(":8080", r)
}
