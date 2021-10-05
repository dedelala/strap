package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	*template.Template
}

func NewServer(t *template.Template) (*Server, error) {
	s := &Server{
		Template: t,
	}
	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.handleGet(rw, req)
	default:
		methods := []string{http.MethodGet}
		rw.Header().Set("Accept", strings.Join(methods, ", "))
		http.Error(rw, "invalid request method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGet(rw http.ResponseWriter, req *http.Request) {
	var bb bytes.Buffer
	contentType := selectContentType(req, "text/plain", "text/html")
	switch contentType {
	case "text/plain":
		fmt.Fprintln(&bb, "plain text response!")
	case "text/html":
		err := s.Execute(&bb, nil)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(rw, "failed to select content type", http.StatusInternalServerError)
		return
	}
	br := bytes.NewReader(bb.Bytes())
	http.ServeContent(rw, req, "", time.Now(), br)
}

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

//go:embed index.html
var embedIndexHTML string

//go:embed static
var embedStatic embed.FS

func main() {
	var (
		listen = ":8080"
	)
	flag.StringVar(&listen, "listen", listen, "server listen address")
	flag.Parse()

	funcs := map[string]interface{}{}
	t, err := template.New("t").Funcs(funcs).Parse(embedIndexHTML)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	sv, err := NewServer(t)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/", http.FileServer(http.FS(embedStatic))))
	mux.Handle("/", sv)

	lgr := NewLogger(mux)

	hsv := &http.Server{
		Addr:    listen,
		Handler: lgr,
	}
	log.Fatal(hsv.ListenAndServe())
}
