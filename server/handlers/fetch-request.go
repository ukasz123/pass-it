package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"pass-it/cache"
	"pass-it/server/models"
	"strings"
	"time"

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
	// 👇👇
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}
	id, _ := uuid.NewRandom()
	idString := id.String()
	w.Header().Set("Content-Type", "text/event-stream")

	idMessage := formatEvent("id", idString)
	w.Write([]byte(idMessage))
	flusher.Flush()

	ch := make(chan struct {
		string
		bool
	})

	go waitForConfirmation(r.Context(), idString, h.payloadMessageChannel, ch)
	timeout := time.AfterFunc(cache.DefaultTTL, func() {
		ch <- struct {
			string
			bool
		}{"Timeout", true}
	})
	defer timeout.Stop()

	for i := range ch {
		if i.bool {
			return
		}
		priceMessage := formatEvent("confirmation", fmt.Sprintf("%s", i.string))
		w.Write([]byte(priceMessage))
		flusher.Flush()
	}
}

func waitForConfirmation(context context.Context, id string, payloadInput chan models.PayloadMessage[string, string], output chan struct {
	string
	bool
}) {
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

func generateEvents(ctx context.Context, priceCh chan<- int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	ticker := time.NewTicker(time.Second)

outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		case <-ticker.C:
			p := r.Intn(100)
			priceCh <- p
		}
	}

	ticker.Stop()

	close(priceCh)
}

func formatEvent(event string, data string) string {
	lines := strings.Split(data, "\n")

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	for _, l := range lines {
		sb.WriteString(fmt.Sprintf("data: %s\n", l))
	}
	sb.WriteString("\n")
	return sb.String()
}
