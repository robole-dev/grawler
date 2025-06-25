package cmd

import (
	"bytes"
	"io"
	"testing"
)

func TestCrawler(t *testing.T) {

	//t.Fatalf("oh no")

	cmd := grawlCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"grawl", "https://robole.de"})
	cmd.Execute()
	out, err := io.ReadAll(b)

	if err != nil {
		t.Fatal(err)
	}

	if string(out) != "hi" {
		t.Fatalf("expected \"%s\" got \"%s\"", "hi", string(out))
	}

	//// Simuliere eine kleine HTML-Seite
	//testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte(`<html><body><a href="/next">Link</a></body></html>`))
	//}))
	//defer testServer.Close()
	//
	//visited := []string{}
	//c := colly.NewCollector()
	//
	//c.OnHTML("a", func(e *colly.HTMLElement) {
	//	visited = append(visited, e.Attr("href"))
	//})
	//
	//err := c.Visit(testServer.URL)
	//if err != nil {
	//	t.Fatalf("Visit failed: %v", err)
	//}
	//
	//if len(visited) != 1 || visited[0] != "/next" {
	//	t.Errorf("Unexpected visited links: %v", visited)
	//}
}
