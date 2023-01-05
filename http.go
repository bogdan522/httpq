package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPQ struct {
	RxBytes  int                      // number of bytes (message body) consumed
	TxBytes  int                      // number of bytes (message body) published
	PubFails int                      // number of publish failures
	SubFails int                      // number of subscribe failures
	queue    map[string]chan []string // message queue
}

func (h *HTTPQ) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/{topic}", h.Consume().ServeHTTP)
	r.Post("/{topic}", h.Publish().ServeHTTP)
	r.Get("/stats", h.Stats().ServeHTTP)

	return r
}

func (h *HTTPQ) Publish() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.queue == nil {
			h.queue = make(map[string]chan []string)
		}

		if r.ContentLength == 0 {
			h.PubFails++
			w.WriteHeader(400)
			return
		}

		msg, err := io.ReadAll(r.Body)
		if len(msg) <= 0 {
			h.PubFails++
			w.WriteHeader(400)
			return
		}

		if err != nil {
			h.PubFails++
			w.WriteHeader(400)
			return
		}

		topic := r.URL.Path

		if _, ok := h.queue[topic]; !ok {
			h.queue[topic] = make(chan []string, 1)
		}

		h.queue[topic] <- []string{string(msg)}
		h.TxBytes += len([]byte(msg))

		w.WriteHeader(201)
	})
}

func (h *HTTPQ) Consume() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Path

		val, ok := h.queue[topic]

		if !ok || len(val) <= 0 {
			h.SubFails++
			w.WriteHeader(400)
			return
		}

		msg := <-h.queue[topic]

		h.RxBytes += len([]byte(msg[0]))

		w.Write([]byte(msg[0]))
	})
}

func (h *HTTPQ) Stats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h)
	})
}
