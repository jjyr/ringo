package ringo

import (
	"bufio"
	"net"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	written    bool
	statusCode int
}

// test interfaces
var _ http.ResponseWriter = &ResponseWriter{}
var _ http.Hijacker = &ResponseWriter{}
var _ http.CloseNotifier = &ResponseWriter{}
var _ http.Flusher = &ResponseWriter{}

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

func (w *ResponseWriter) Write(content []byte) (n int, err error) {
	w.written = true
	n, err = w.ResponseWriter.Write(content)
	return
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *ResponseWriter) writeHeaderNow() {
	if !w.Written() {
		w.written = true
		w.ResponseWriter.WriteHeader(w.statusCode)
	}
}

func (w *ResponseWriter) Written() bool {
	return w.written
}

// Implements the http.Hijack interface
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// Implements the http.CloseNotify interface
func (w *ResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Implements the http.Flush interface
func (w *ResponseWriter) Flush() {
	if w.Written() {
		w.ResponseWriter.(http.Flusher).Flush()
	}
}
