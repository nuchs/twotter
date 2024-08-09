package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func NewServer(ctx context.Context, twotter *Twotter) (*http.Server, error) {
	twotterHandler, err := TwotterHandler(twotter)
	if err != nil {
		return nil, fmt.Errorf("Failed to create twotter handler: %w", err)
	}
	userHandler, err := UserHandler(twotter)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user handler: %w", err)
	}
	adminHandler, err := AdminHandler(twotter)
	if err != nil {
		return nil, fmt.Errorf("Failed to create admin handler: %w", err)
	}

	mws := ChainMiddleWare(loggerMW, makeStaticMW("static"))
	mux := http.NewServeMux()
	mux.Handle("/", mws(twotterHandler))
	mux.Handle("/user", mws(userHandler))
	mux.Handle("/admin", mws(adminHandler))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}, nil
}

func loggerMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func makeStaticMW(dir string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/")
			path = strings.TrimPrefix(path, ".")
			if path == "" {
				h.ServeHTTP(w, r)
				return
			}

			path = fmt.Sprintf("%s/%s", dir, path)
			file, err := os.Open(path)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			defer file.Close()

			log.Printf("Serving file: %s", path)
			if _, err = io.Copy(w, file); err != nil {
				log.Printf("Error copying: %v", err)
			}
		})
	}
}
