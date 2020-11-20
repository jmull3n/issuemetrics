package server

import (
	"errors"
	"bufio"
	"net"
	"net/http"
)

// ResponseWriter wraps the standard http.ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

// NewResponseWriter returns ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// Status provides an easy way to retrieve the status code
func (w *ResponseWriter) Status() int {
	return w.status
}


// Write satisfies the http.ResponseWriter interface
func (w *ResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	return size, err
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows catching the status code
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
// Hijack implements the hijack interaface
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    h, ok := w.ResponseWriter.(http.Hijacker)
    if !ok {
        return nil, nil, errors.New("hijack not supported")
    }
    return h.Hijack()
}