package main

import (
	"compress/gzip"
	"env-up-app/backend/repository"
	"env-up-app/backend/webservices"
	"log"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jamesrr39/goutil/httpextra"
	"github.com/jamesrr39/goutil/userextra"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	configLocation := kingpin.Flag("config-location", "filepath to where the configuration is stored").Default("~/.config/github.com/jamesrr39/env-up-main/default.yaml").String()
	kingpin.Parse()

	expandedConfigLocation, err := userextra.ExpandUser(*configLocation)
	if err != nil {
		log.Fatalf("failed to expand '%s'. Error: '%s'\n", *configLocation, err)
	}

	err = os.MkdirAll(filepath.Dir(expandedConfigLocation), 0700)
	if err != nil {
		log.Fatalf("failed to create '%s'. Error: '%s'\n", expandedConfigLocation, err)
	}

	configRepository, err := repository.NewConfigRepository(expandedConfigLocation)
	if err != nil {
		log.Fatalf("failed to read configuration from '%s'. Error: '%s'\n", expandedConfigLocation, err)
	}

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Compress(gzip.DefaultCompression))
	mainRouter.Route("/api/", func(r chi.Router) {
		r.Mount("/config", webservices.NewConfigWebService(configRepository))
		r.Mount("/environment", webservices.NewEnvironmentWebService(repository.NewEnvironmentRepository()))
	})
	mainRouter.Mount("/", webservices.NewStaticAssetsHandler())

	server := httpextra.NewServerWithTimeouts()
	server.Addr = "localhost:9010"
	server.Handler = mainRouter
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
