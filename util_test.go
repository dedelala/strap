package main

import (
	"net/http"
	"testing"
)

func TestSelectContentType(t *testing.T) {
	mediaTypes := []string{"text/plain", "text/html"}
	var ss = []struct {
		name   string
		accept []string
		exp    string
	}{
		{
			name: "chrome header select html",
			accept: []string{
				"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			},
			exp: "text/html",
		},
		{
			name: "no header select plain",
			exp:  "text/plain",
		},
	}
	for _, s := range ss {
		t.Run(s.name, func(t *testing.T) {
			req := &http.Request{
				Header: http.Header{"Accept": s.accept},
			}
			res := selectContentType(req, mediaTypes...)
			if s.exp != res {
				t.Errorf("expected %q got %q", s.exp, res)
			}
		})
	}
}
