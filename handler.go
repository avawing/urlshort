package urlshort

import (
	"fmt"
	yamlv2 "gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//type that is essentially a function
	// we do not have to cast
	// it's return a handler func
	return func(w http.ResponseWriter, r *http.Request) {
		// extract path from url
		path := r.URL.Path
		// if we can match a path => redirect
		// sugar: if you can find a key in the map, it will be true
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			//stop function
			return
		}
		// otherwise call a fallback
		fallback.ServeHTTP(w, r)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse yaml
	pathURLS, err := parseYAML(yml)
	if err != nil {
		fmt.Println("Oops")
	}
	// convert yaml array into map
	pathsToUrls := buildMap(pathURLS)
	// return map handler

	// unmarshall is turning data from yaml => something else
	// marshall converts data => yaml

	return MapHandler(pathsToUrls, fallback), nil
}

type pathURL struct {
	path string `yaml:"path"`
	url  string `yaml:"url"`
}

func buildMap(pathurls []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	// don't need index
	for _, pu := range pathURLS {
		pathsToUrls[pu.path] = pu.url

	}
	return pathsToUrls
}

func parseYAML(data []byte) ([]pathURL, error) {
	var pathURLS []pathURL
	err := yamlv2.Unmarshal(data, &pathURLS)
	if err != nil {
		return nil, err
	}
	return pathURLS, nil
}
