package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPQ struct {
	RxBytes  int                 // number of bytes (message body) consumed
	TxBytes  int                 // number of bytes (message body) published
	PubFails int                 // number of publish failures
	SubFails int                 // number of subscribe failures
	queue    map[string][]string // message queue
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
			h.queue = make(map[string][]string)
		}

		if r.ContentLength == 0 {
			h.PubFails++
			w.WriteHeader(400)
			return
		}

		msg, err := io.ReadAll(r.Body)

		if err != nil {
			h.PubFails++
			w.WriteHeader(400)
			return
		}

		topic := r.URL.Path

		val, ok := h.queue[topic]

		if ok {
			val = append(val, string(msg))
			h.queue[topic] = val
		} else {
			h.queue[topic] = append(make([]string, 0), string(msg))
		}

		h.TxBytes += len([]byte(msg))
		w.Write(msg)
		w.WriteHeader(200)

	})
}

func (h *HTTPQ) Consume() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Path

		val, ok := h.queue[topic]

		if !ok || len(val) <= 0 {
			h.SubFails++
			w.WriteHeader(400)
			w.Write([]byte("Array empty"))
			return
		}

		msg := val[0]
		h.queue[topic] = val[1:]
		h.RxBytes += len([]byte(msg))

		w.Write([]byte(msg))
		w.WriteHeader(200)

	})
}

func (h *HTTPQ) Stats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h)
	})
}
