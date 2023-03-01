package handlers

import (
	"io"
	"net/http"
	"time"

	"playstoreproxy/internal/log"
)

type PlayStoreHandler struct{}

// NewNotFoundHandler returns PlayStoreHandler interactor
func NewPlayStoreHandler() *PlayStoreHandler {
	return &PlayStoreHandler{}
}

func (h *PlayStoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://play.google.com/store/apps/details?id=ru.blizko_mobile")
	if err != nil {
		log.Errorf("Failed to load play store: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Errorf("Failed to load play store: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
