package handlers

import (
	"hiv_mind/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	serviceProvider := app.NewServiceProvider()
	metricHandler, err := NewMetricHandler(serviceProvider, "../../templates")
	require.NoError(t, err, "Не удалось создать хендлер")

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", metricHandler.Update)

	srv := httptest.NewServer(r)
	defer srv.Close()

	data := []struct {
		tName        string
		method       string
		path         string
		contentType  string
		responseCode int
	}{
		{
			tName:        "Valid_POST_Request",
			method:       http.MethodPost,
			path:         "/update/gauge/GCSys/1.85424e+06",
			contentType:  "text/plain",
			responseCode: http.StatusOK,
		},
		{
			tName:        "Request_With_NOT_ALLOWDED_METHOD",
			method:       http.MethodGet,
			path:         "/update/gauge/GCSys/1.85424e+06",
			contentType:  "text/plain",
			responseCode: http.StatusMethodNotAllowed,
		},
		{
			tName:        "BAD_REQUEST_WITHOUT_VALUE",
			method:       http.MethodPost,
			path:         "/update/gauge/GCSys/",
			contentType:  "text/plain",
			responseCode: http.StatusNotFound,
		},
		{
			tName:        "BAD_REQUEST_WITH_NOT_ALLOWDED_METRIC_TYPE",
			method:       http.MethodPost,
			path:         "/update/notallowded/GCSys/123",
			contentType:  "text/plain",
			responseCode: http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.tName, func(t *testing.T) {
			fullURL := srv.URL + d.path

			client := resty.New().R()
			client.Method = d.method
			client.URL = fullURL
			response, err := client.SetHeader("Content-Type", d.contentType).Send()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, d.responseCode, response.StatusCode())
		})
	}
}
