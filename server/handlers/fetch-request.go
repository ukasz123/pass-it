package handlers

import (
	"context"
	"log"
	"net/http"
	"pass-it/cache"
	"pass-it/server/models"
	"pass-it/server/templates"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type FetchRequestHandler struct {
	cache                 cache.Cache[models.DefaultStoredData]
	payloadMessageChannel chan models.PayloadMessage[string, string]
}

func NewFetchRequestHandler(c cache.Cache[models.DefaultStoredData], p chan models.PayloadMessage[string, string]) *FetchRequestHandler {
	return &FetchRequestHandler{cache: c, payloadMessageChannel: p}
}

func (h *FetchRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer log.Default().Println("Fetch request done")
	// 👇👇
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}
	var rend = componentRenderer{
		w: w,
		f: flusher,
		c: r.Context(),
	}

	id, _ := uuid.NewRandom()
	idString := id.String()
	w.Header().Set("Content-Type", "text/event-stream")

	log.Default().Printf("SSE ID: %s", idString)
	if err := rend.renderComponentToSSE(templates.SessionCode(idString)); err != nil {
		log.Default().Printf("Error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ch := make(chan struct {
		string
		bool
	})

	timeout := time.AfterFunc(cache.DefaultTTL, func() {
		rend.renderComponentToSSE(templates.Timeout())
		log.Default().Println("Timeout")
		ch <- struct {
			string
			bool
		}{"", true}
	})
	defer timeout.Stop()
	
	go waitForConfirmation(r.Context(), idString, h.payloadMessageChannel, ch)

	for i := range ch {
		log.Default().Printf("Confirmation: %s, %v\n", i.string, i.bool)
		if i.bool {
			return
		}
		rend.renderComponentToSSE(templates.Secret(i.string))
	}
}

func waitForConfirmation(context context.Context, id string, payloadInput chan models.PayloadMessage[string, string], output chan struct {
	string
	bool
}) {
	defer log.Default().Println("Confirmation request done")
	defer close(output)
	for {
		select {
		case <-context.Done():
			return
		case msg := <-payloadInput:
			if msg.Addr == id {
				output <- struct {
					string
					bool
				}{msg.Payload, false}
				return
			}
		}
	}
}

type componentRenderer struct {
	w http.ResponseWriter
	f http.Flusher
	c context.Context
}

func (r *componentRenderer) renderComponentToSSE(c templ.Component) error {
	r.w.Write([]byte("data:"))

	defer r.f.Flush()
	defer r.w.Write([]byte("\n\n"))

	return c.Render(r.c, r.w)
}
