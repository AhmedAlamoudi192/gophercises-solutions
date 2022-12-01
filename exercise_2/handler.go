package urlshort

import (
	"fmt"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		v, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, v, 301)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type yamlScheme struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

//	pathsToUrls := map[string]string{
//		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
//		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
//	}
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var out []yamlScheme
	err := yaml.Unmarshal(yml, &out)
	if err != nil {
		fmt.Println(err)
		return func(w http.ResponseWriter, r *http.Request) { fallback.ServeHTTP(w, r) }, err
	}
	pathsToUrls := map[string]string{}
	for _, i := range out {
		pathsToUrls[i.Path] = i.Url
	}
	return MapHandler(pathsToUrls, fallback), nil

}
