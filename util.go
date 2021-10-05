package main

import (
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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
