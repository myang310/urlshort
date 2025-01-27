package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	json "encoding/json"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		path, ok := pathsToUrls[reader.URL.Path]
		if ok {
			http.Redirect(writer,reader,path, http.StatusPermanentRedirect)
		} else {
			print("Using fallback handler...")
			fallback.ServeHTTP(writer, reader)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(data []byte) (t []RedirectionMap, err error) {
	err = yaml.Unmarshal(data, &t)

	return t, err
}

func buildMap(parsedYAML []RedirectionMap) map[string]string {
	m := make(map[string]string)
	for _, entry := range parsedYAML {
		m[entry.Path] = entry.Url
	}

	return m
}

type RedirectionMap struct {
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error)  {
	parsedJSON, err := parseJSON(data)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(data []byte) (t []RedirectionMap, err error) {
	err = json.Unmarshal(data, &t)

	return t, err
}