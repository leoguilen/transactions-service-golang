package handlers

import (
	"bytes"
	"net/http"
)

type ResponseInterceptor struct {
	http.ResponseWriter
	statusCode    int
	bodyBuffer    bytes.Buffer
	bodyCaptured  int64
	headerWritten bool
}

func NewResponseInterceptor(w http.ResponseWriter) *ResponseInterceptor {
	return &ResponseInterceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (ri *ResponseInterceptor) WriteHeader(code int) {
	ri.statusCode = code
	ri.headerWritten = true
	ri.ResponseWriter.WriteHeader(code)
}

func (ri *ResponseInterceptor) Write(b []byte) (int, error) {
	return ri.ResponseWriter.Write(b)
}

func (ri *ResponseInterceptor) GetStatusCode() int {
	return ri.statusCode
}
