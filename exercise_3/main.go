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

func registerRoute(mux *http.ServeMux, arc string, stories map[string]StoryArc) {
	templ:=`<body style="text-align:center"><h1>{{.Title}}</h1>{{range .Story}}<p>{{.}}</p>{{end}}<div style="display:flex;gap:2rem;justify-content:center">{{range .Options}}<a href="{{.Arc}}">{{.Text}}</a>{{end}}</div></body>`
	tmpl, err := template.New(arc).Parse(templ)
	checknil(err)
	if strings.Compare(arc, "intro")!=0 {	
		mux.HandleFunc(fmt.Sprintf("/%s",arc), func(w http.ResponseWriter, h *http.Request) {
			tmpl.Execute(w, stories[arc])
		})
	}else {
		mux.HandleFunc("/", func(w http.ResponseWriter, h *http.Request) {
			tmpl.Execute(w, stories[arc])
		})
	}
}

func httpHandler() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func main() {
	jsonBlob, os_err := os.ReadFile("./gopher.json")
	checknil(os_err)
	var stories map[string]StoryArc
	jsonErr := json.Unmarshal(jsonBlob, &stories)
	checknil(jsonErr)
	mux := httpHandler()
	for k := range stories { 
		registerRoute(mux,k,stories)
	}
	http.ListenAndServe("localhost:8080",mux)
}
