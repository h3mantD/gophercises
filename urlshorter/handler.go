package urlshorter

import (
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func MapHandler(urls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, err := urls[r.URL.Path]; err {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(fallback http.Handler) (http.HandlerFunc, error) {
	f, err := os.ReadFile("./routes/urls.yaml")
	if err != nil {
		return nil, err
	}

	var raw []map[string]string
	if err := yaml.Unmarshal(f, &raw); err != nil {
		return nil, err
	}

	pathsToUrls := map[string]string{}
	for _, m := range raw {
		pathsToUrls[m["path"]] = m["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(fallback http.Handler) (http.HandlerFunc, error) {
	f, err := os.ReadFile("./routes/urls.json")
	if err != nil {
		return nil, err
	}

	raw := []map[string]string{}
	if err := json.Unmarshal(f, &raw); err != nil {
		return nil, err
	}

	pathsToUrls := map[string]string{}
	for _, m := range raw {
		pathsToUrls[m["path"]] = m["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}
