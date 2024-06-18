package app

import (
	"crypto/md5"
	"encoding/hex"
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

	hash := md5.New()
	hash.Write([]byte(urlValue))
	urlID := hex.EncodeToString(hash.Sum(nil)[:8])
	urls[urlID] = urlValue

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://" + req.Host + "/" + urlID))
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
