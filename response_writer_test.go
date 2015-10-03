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
	w.content = content
	return len(content), nil
}

func (w *FakeResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
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
	w.CleanBuffer()
	if w.Written() || w.content != nil || w.statusCode != 0 {
		t.Errorf("Content and written state should cleared")
	}

	if w.Flushed() {
		t.Errorf("Should not be flushed")
	}

	w.WriteHeader(404)
	if !w.Written() || w.Flushed() || w.statusCode != 404 || fw.statusCode != 0 {
		t.Errorf("Value not correct")
	}
	w.Write([]byte("OK"))
	w.Flush()
	if !w.Written() || !w.Flushed() || w.statusCode != 404 || fw.statusCode != 404 || string(fw.content) != "OK" {
		t.Errorf("Value not correct")
	}

	w.CleanBuffer()
	if w.Written() || !w.Flushed() || w.statusCode != 0 || w.content != nil {
		t.Errorf("Clear buffer not work, %+v", *w)
	}
}
