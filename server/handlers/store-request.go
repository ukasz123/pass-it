package handlers

import (
	"log"
	"net/http"
	"pass-it/cache"
	"pass-it/crypto"
	"pass-it/server/models"

	"github.com/gorilla/mux"
)

type StoreRequestHandler struct {
	cache cache.Cache[models.DefaultStoredData]
}

func NewStoreRequestHandler(c cache.Cache[models.DefaultStoredData]) *StoreRequestHandler {
	return &StoreRequestHandler{cache: c}
}

func (h *StoreRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if h.cache.Contains(id) {
		w.WriteHeader(http.StatusConflict)
		log.Printf("Key %s already exists\n", id)
		return
	}
	var key = r.FormValue("key")
	var payload = r.FormValue("payload")
	if key == "" || payload == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("`key` or `payload` is empty"))
		return
	}
	var decodedKey, err = crypto.DecodePublicKey(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid `key`"))
	}
	h.cache.Set(id, models.DefaultStoredData{Key: decodedKey, Payload: payload})
	w.WriteHeader(http.StatusCreated)
}
