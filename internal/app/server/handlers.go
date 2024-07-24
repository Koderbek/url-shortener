package server

import (
	"encoding/json"
	"github.com/Koderbek/url-shortener/internal/app/config"
	"github.com/Koderbek/url-shortener/internal/app/models"
	"github.com/Koderbek/url-shortener/internal/app/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
)

var urls = make(map[string]string)

func shorten(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	urlValue := string(body)
	_, err = url.ParseRequestURI(urlValue)
	if err != nil {
		http.Error(res, "Bad URL value", http.StatusBadRequest)
		return
	}

	urlID := utils.HashGenerator(urlValue, 8)
	urls[urlID] = urlValue

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Config.Flags.ShortenedAddress + "/" + urlID))
}

func apiShorten(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	var shortenRequest models.ShortenRequest
	if err := json.NewDecoder(req.Body).Decode(&shortenRequest); err != nil {
		http.Error(res, "Decode request error", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(shortenRequest.URL); err != nil {
		http.Error(res, "Bad URL value", http.StatusBadRequest)
		return
	}

	urlID := utils.HashGenerator(shortenRequest.URL, 8)
	urls[urlID] = shortenRequest.URL

	shortenResponse := models.ShortenResponse{Result: config.Config.Flags.ShortenedAddress + "/" + urlID}
	jsonResponse, err := json.Marshal(shortenResponse)
	if err != nil {
		http.Error(res, "Encode response error", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResponse)
}

func findURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	urlID := chi.URLParam(req, "id")
	urlValue := urls[urlID]
	if urlValue == "" {
		http.Error(res, "Url not found", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", urlValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
