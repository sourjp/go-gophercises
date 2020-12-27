package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToURLs map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := pathsToURLs[r.URL.Path]
		if len(url) > 0 {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
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
	var pathsToURLsList []PathToURL
	if err := yaml.Unmarshal(yml, &pathsToURLsList); err != nil {
		return nil, err
	}
	pathsToURLs := buildMapper(pathsToURLsList)
	return MapHandler(pathsToURLs, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToURLList []PathToURL
	if err := json.Unmarshal(jsn, &pathsToURLList); err != nil {
		return nil, err
	}
	pathsToURLs := buildMapper(pathsToURLList)
	return MapHandler(pathsToURLs, fallback), nil
}

// DBHandler will get Paths to URLs from DB and then return
// an http.HandlerFunc
func DBHandler(conn Conn, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToURLList, err := conn.GetPathsToURLsByDB()
	if err != nil {
		return nil, err
	}
	pathsToURLs := buildMapper(pathsToURLList)
	return MapHandler(pathsToURLs, fallback), nil
}

// PathToURL mapping short URL and target URL of struct
type PathToURL struct {
	URL  string `yaml:"url" db:"url"`
	Path string `yaml:"path" db:"path"`
}

func buildMapper(pathToURLs []PathToURL) map[string]string {
	mapper := map[string]string{}
	for _, pathToURL := range pathToURLs {
		mapper[pathToURL.Path] = pathToURL.URL
	}
	return mapper
}
