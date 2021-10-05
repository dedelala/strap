package main

import (
	"log"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Logger struct {
	http.Handler
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	lrw := &logResponseWriter{
		ResponseWriter: rw,
		status:         http.StatusOK,
	}
	l.Handler.ServeHTTP(lrw, req)
	log.Printf("%s %s %s %d %d %s", req.RemoteAddr, req.Method, req.URL.Path,
		lrw.status, lrw.length, lrw.Header().Get("Content-Type"))
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
	length int64
}

func (rw *logResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.length += int64(n)
	return n, err
}

func (rw *logResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func selectContentType(req *http.Request, mediaTypes ...string) string {
	var (
		acceptValues = req.Header.Values("Accept")
		acceptTypes  []string
	)
	for _, v := range acceptValues {
		acceptTypes = append(acceptTypes, strings.Split(v, ",")...)
	}
	sort.SliceStable(acceptTypes, func(i, j int) bool {
		qi := parseMediaTypeQ(acceptTypes[i])
		qj := parseMediaTypeQ(acceptTypes[j])
		return qi > qj
	})
	for _, acceptType := range acceptTypes {
		for _, mediaType := range mediaTypes {
			if matchContentType(acceptType, mediaType) {
				return mediaType
			}
		}
	}
	if len(mediaTypes) == 0 {
		return ""
	}
	return mediaTypes[0]
}

func matchContentType(acceptValue, mediaType string) bool {
	acceptType, _, err := mime.ParseMediaType(acceptValue)
	switch {
	case err != nil:
		return false
	case mediaType == acceptType:
		return true
	case acceptType == "*/*":
		return true
	case !strings.HasSuffix(acceptType, "/*"):
		return false
	default:
		acceptType = strings.TrimSuffix(acceptType, "*")
		return strings.HasPrefix(mediaType, acceptType)
	}
}

func parseMediaTypeQ(acceptValue string) float64 {
	_, params, _ := mime.ParseMediaType(acceptValue)
	qval, ok := params["q"]
	if !ok {
		return 1.0
	}
	q, err := strconv.ParseFloat(qval, 64)
	if err != nil {
		return 1.0
	}
	return q
}
