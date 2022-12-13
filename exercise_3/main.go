package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func checknil(e error) {
	if e != nil {
		fmt.Println("error:", e)
	}
}

func httpHandler() *http.ServeMux {
	mux := http.NewServeMux()
	tmpl, err := template.New("Entry").Parse(`<h1>hello world</h1>`)
	checknil(err)
	mux.HandleFunc("/", func(w http.ResponseWriter, h *http.Request) {
		tmpl.Execute(w, nil)
	})
	return mux
}

func main() {
	jsonBlob, os_err := os.ReadFile("./gopher.json")
	checknil(os_err)
	var stories map[string]StoryArc
	json_err := json.Unmarshal(jsonBlob, &stories)
	checknil(json_err)
	fmt.Printf("%s\n", stories["intro"].Title)
	fmt.Printf("%s\n", strings.Join(stories["intro"].Story, ""))
	fmt.Printf("- %s\n", stories["intro"].Options[0].Text)
	fmt.Printf("- %s\n", stories["intro"].Options[1].Text)
}
