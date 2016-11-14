package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		url      string
		handler  http.HandlerFunc
		expected int
	}{
		{
			"/",
			IndexHandler,
			200,
		},
		{
			"/index",
			IndexHandler,
			200,
		},
		{
			"/setting",
			SettingHandler,
			200,
		},
		{
			"/log",
			LogHandler,
			200,
		},
		{
			"/robots.txt",
			RobotsHandler,
			200,
		},
	}

	rr := httptest.NewRecorder()
	for i, test := range tests {
		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Errorf("test #%d, expected %d, get %v", i, test.expected, err)
		}
		handler := http.HandlerFunc(test.handler)
		handler.ServeHTTP(rr, req)
		if rr.Code != test.expected {
			t.Errorf("test #%d, expected %d, get %d", i, test.expected, rr.Code)
		}
	}
}
