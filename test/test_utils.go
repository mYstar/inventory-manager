package test

import (
	"inventory_manager/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
)

func sendTestRequest(method, url, body string) *httptest.ResponseRecorder {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	router.ServeHTTP(w, req)

	return w
}
