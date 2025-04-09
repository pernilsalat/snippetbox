package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/platform/web"
	"snippetbox/pkg/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	web.SecureHeaders(next).ServeHTTP(recorder, r)
	rsp := recorder.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rsp.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rsp.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	assert.Equal(t, rsp.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	assert.Equal(t, rsp.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	assert.Equal(t, rsp.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, http.StatusOK, rsp.StatusCode)

	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK", string(body))
}
