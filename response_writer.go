package ringo

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	written bool
}

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, written: false}
}

func (w *ResponseWriter) Write(content []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(content)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.written = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Written() bool {
	return w.written
}
