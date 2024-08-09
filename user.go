package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func UserHandler(twotter *Twotter) (http.Handler, error) {
	tmpl := template.New("layout.html")
	_, err := tmpl.ParseFiles(
		"views/layout.html",
		"views/users.html",
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse templates: %w", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, twotter); err != nil {
			log.Printf("Error executing template: %v", err)
		}
	}), nil
}
