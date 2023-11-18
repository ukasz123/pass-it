package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"pass-it/cache"
	"pass-it/crypto"
	"pass-it/server/models"

	"github.com/gorilla/mux"
)

type ConfirmRequestHandler struct {
	cache               cache.Cache[models.DefaultStoredData]
	confirmationChannel chan models.PayloadMessage[string, string]
}

func NewConfirmRequestHandler(cache cache.Cache[models.DefaultStoredData],
	confirmationChannel chan models.PayloadMessage[string, string]) *ConfirmRequestHandler {
	return &ConfirmRequestHandler{
		cache:               cache,
		confirmationChannel: confirmationChannel,
	}
}

func (h *ConfirmRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !h.cache.Contains(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var signature = r.FormValue("signature")
	var sessionId = r.FormValue("session_id")
	if signature == "" || sessionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("`signature` or `session_id` is empty"))
		return
	}

	var storedData = h.cache.Get(id)
	var key = storedData.Value().Key
	signatureDecoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = crypto.CheckSignature(key, signatureDecoded, []byte(sessionId))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.confirmationChannel <- models.PayloadMessage[string, string]{Addr: sessionId, Payload: storedData.Value().Payload}
	w.WriteHeader(http.StatusFound)
}
