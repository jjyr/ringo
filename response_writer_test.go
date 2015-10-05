package ringo

import (
	"net/http"
	"testing"
)

type FakeResponseWriter struct {
	statusCode int
	content    []byte
	http.ResponseWriter
}

func (w *FakeResponseWriter) Write(content []byte) (int, error) {
	return len(content), nil
}

func (w *FakeResponseWriter) WriteHeader(statusCode int) {
}

func TestResponseWriter(t *testing.T) {
	fw := FakeResponseWriter{}
	w := newResponseWriter(&fw)
	if w.Written() {
		t.Errorf("Should not written")
	}
	w.Write([]byte("hello"))
	if !w.Written() {
		t.Errorf("Should written")
	}

	w = newResponseWriter(&fw)
	w.WriteHeader(404)
	if w.Written() {
		t.Errorf("Should not written")
	}
	w.writeHeaderNow()

	if !w.Written() {
		t.Errorf("Should written")
	}
}
