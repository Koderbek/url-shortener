package server

import (
	"github.com/Koderbek/url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return ts.Client().Do(req)
}

func Test_shorten(t *testing.T) {
	ts := httptest.NewServer(newServer().router)
	defer ts.Close()

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
				response:    config.Config.Flags.ShortenedAddress + "/d7115cf9972dcaf2",
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
			resp, err := testRequest(t, ts, tt.method, "/", strings.NewReader(tt.body))
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, resp.StatusCode, tt.want.code)
			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, string(respBody), tt.want.response)
		})
	}
}

func Test_apiShorten(t *testing.T) {
	ts := httptest.NewServer(newServer().router)
	defer ts.Close()

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
			body:   `{"url": "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr"}`,
			want: want{
				code:        http.StatusCreated,
				response:    `{"result":"` + config.Config.Flags.ShortenedAddress + `/d7115cf9972dcaf2"}`,
				contentType: "application/json",
			},
		},
		{
			name:   "negative test #1",
			method: http.MethodGet,
			body:   `{"url": "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr"}`,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodPost,
			body:   "http://kgeus60l.com/avlpcuyp2iq/dphj1mszqiqvi/bp9sfaxr",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Decode request error\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #3",
			method: http.MethodPost,
			body:   `{"url": "1"}`,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad URL value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testRequest(t, ts, tt.method, "/api/shorten", strings.NewReader(tt.body))
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, resp.StatusCode, tt.want.code)
			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, string(respBody), tt.want.response)
		})
	}
}

func Test_findURL(t *testing.T) {
	ts := httptest.NewServer(newServer().router)
	defer ts.Close()

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

			resp, err := testRequest(t, ts, tt.method, "/"+tt.id, nil)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, resp.StatusCode, tt.want.code)
			if tt.want.isError {
				assert.Equal(t, string(respBody), tt.want.response)
			} else {
				assert.Equal(t, resp.Header.Get("Location"), tt.want.response)
			}
		})
	}
}
