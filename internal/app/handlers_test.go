package app

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_shorten(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPost,
			body:   "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr",
			want: want{
				code:        http.StatusCreated,
				response:    "http://example.com/d7115cf9972dcaf2",
				contentType: "text/plain",
			},
		},
		{
			name:   "negative test #1",
			method: http.MethodGet,
			body:   "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodPost,
			body:   "12233",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad URL value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))

			// создаём новый Recorder
			w := httptest.NewRecorder()
			w.Header().Set("Content-Type", "text/plain")

			shorten(w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, string(resBody), tt.want.response)
		})
	}
}

func Test_findURL(t *testing.T) {
	type want struct {
		isError  bool
		code     int
		response string
	}

	tests := []struct {
		name   string
		method string
		id     string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			id:     "d7115cf9972dcaf2",
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr",
			},
		},
		{
			name:   "negative test #1",
			method: http.MethodPost,
			id:     "d7115cf9972dcaf2",
			want: want{
				isError:  true,
				code:     http.StatusBadRequest,
				response: "Only GET requests are allowed!\n",
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodGet,
			id:     "12233",
			want: want{
				isError:  true,
				code:     http.StatusBadRequest,
				response: "Url not found\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.want.isError {
				urls[tt.id] = tt.want.response
			}

			request := httptest.NewRequest(tt.method, "/"+tt.id, nil)
			w := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/{id}", findURL)
			r.ServeHTTP(w, request)

			// Получаем результат
			res := w.Result()
			assert.Equal(t, res.StatusCode, tt.want.code)
			if tt.want.isError {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, string(resBody), tt.want.response)
			} else {
				assert.Equal(t, res.Header.Get("Location"), tt.want.response)
			}
		})
	}
}
