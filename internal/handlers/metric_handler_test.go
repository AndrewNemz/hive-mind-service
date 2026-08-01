package handlers

import (
	"hiv_mind/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	serviceProvider := app.NewServiceProvider()
	data := []struct {
		tName        string
		method       string
		url          string
		contentType  string
		responseCode int
	}{
		{
			tName:        "Valid_POST_Request",
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/gauge/GCSys/1.85424e+06",
			contentType:  "text/plain",
			responseCode: http.StatusOK,
		},
		{
			tName:        "Request_With_NOT_ALLOWDED_METHOD",
			method:       http.MethodGet,
			url:          "http://localhost:8080/update/gauge/GCSys/1.85424e+06",
			contentType:  "text/plain",
			responseCode: http.StatusMethodNotAllowed,
		},
		{
			tName:        "BAD_REQUEST_WITHOUT_VALUE",
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/gauge/GCSys/",
			contentType:  "text/plain",
			responseCode: http.StatusBadRequest,
		},
		{
			tName:        "BAD_REQUEST_WITH_NOT_ALLOWDED_METRIC_TYPE",
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/notallowded/GCSys/123",
			contentType:  "text/plain",
			responseCode: http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.tName, func(t *testing.T) {
			request := httptest.NewRequest(d.method, d.url, nil)
			request.Header.Set("Content-Type", d.contentType)
			writer := httptest.NewRecorder()

			handler := NewMetricHandler(serviceProvider)
			handler.Update(writer, request)

			respone := writer.Result()
			defer respone.Body.Close()

			assert.Equal(t, d.responseCode, respone.StatusCode)
		})
	}
}
