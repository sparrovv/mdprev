package mdprev

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	mdPrev := buildTestMdPrev("#content")
	iHandler := mdFileHandler(mdPrev)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	iHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}
