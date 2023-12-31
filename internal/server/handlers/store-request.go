package handlers

import (
	"log"
	"net/http"
	"pass-it/internal/cache"
	"pass-it/internal/crypto"
	"pass-it/internal/server/models"

	"github.com/gorilla/mux"
)

type StoreRequestHandler struct {
	cache cache.Cache[models.DefaultStoredData]
}

func NewStoreRequestHandler(c cache.Cache[models.DefaultStoredData]) http.Handler {
	return &StoreRequestHandler{cache: c}
}

func (h *StoreRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Processing `%s`", id)
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
		log.Println("`key` or `payload` is empty")
		return
	}
	var decodedKey, err = crypto.DecodePublicKey(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid `key`"))
		log.Println("Invalid `key`")
	}
	h.cache.Set(id, models.DefaultStoredData{Key: decodedKey, Payload: payload})
	log.Default().Printf("Key %s stored\n", id)
	w.WriteHeader(http.StatusCreated)
}
