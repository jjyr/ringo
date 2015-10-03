package ringo

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	written    bool
	statusCode int
	content    []byte
	flushed    bool
}

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

func (w *ResponseWriter) Write(content []byte) (int, error) {
	w.written = true
	w.content = content
	return len(content), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.written = true
	w.statusCode = statusCode
}

func (w *ResponseWriter) Written() bool {
	return w.written
}

func (w *ResponseWriter) Flushed() bool {
	return w.flushed
}

func (w *ResponseWriter) CleanBuffer() {
	w.written = false
	w.statusCode = 0
	w.content = nil
}

func (w *ResponseWriter) Flush() bool {
	flushed := false
	if w.Written() {
		if w.statusCode > 0 {
			w.ResponseWriter.WriteHeader(w.statusCode)
			w.statusCode = 0
			flushed = true
		}

		if len(w.content) > 0 {
			w.ResponseWriter.Write(w.content)
			w.content = nil
			flushed = true
		}
	}

	if flushed {
		w.flushed = true
	}
	return flushed
}
