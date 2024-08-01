package main

import (
	"net/http"

	"github.com/h3mantd/gophercises/urlshorter"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	pathsToUrls := map[string]string{
		"/one": "https://google.com",
	}

	mapHandler := urlshorter.MapHandler(pathsToUrls, mux)

	yamlHandler, err := urlshorter.YAMLHandler(mapHandler)
	if err != nil {
		yamlHandler = mapHandler
	}

	jsonHandler, err := urlshorter.JSONHandler(yamlHandler)
	if err != nil {
		jsonHandler = yamlHandler
	}

	http.ListenAndServe(":8080", jsonHandler)
}
