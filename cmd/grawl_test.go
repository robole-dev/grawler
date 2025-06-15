package cmd

import (
	"github.com/gocolly/colly/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawler(t *testing.T) {
	// Simuliere eine kleine HTML-Seite
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><a href="/next">Link</a></body></html>`))
	}))
	defer testServer.Close()

	visited := []string{}
	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		visited = append(visited, e.Attr("href"))
	})

	err := c.Visit(testServer.URL)
	if err != nil {
		t.Fatalf("Visit failed: %v", err)
	}

	if len(visited) != 1 || visited[0] != "/next" {
		t.Errorf("Unexpected visited links: %v", visited)
	}
}
