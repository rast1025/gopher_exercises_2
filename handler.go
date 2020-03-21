package urlshort

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if destURL, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, destURL, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type Redirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url"  json:"url"`
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil

}

func buildMap(data []Redirect) map[string]string {
	result := make(map[string]string, len(data))
	for _, el := range data {
		result[el.Path] = el.URL
	}
	return result
}

func parseYAML(yamlBytes []byte) ([]Redirect, error) {
	redirect := []Redirect{}
	err := yaml.Unmarshal(yamlBytes, &redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}

func parseJSON(data []byte) ([]Redirect, error) {
	redirect := []Redirect{}
	err := json.Unmarshal(data, &redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseJSON(data)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}
