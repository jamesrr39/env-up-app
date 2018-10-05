package main

import (
	"compress/gzip"
	"env-up-app/backend/repository"
	"env-up-app/backend/webservices"
	"env-up-app/backend/webservices/staticassets"
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jamesrr39/goutil/httpextra"
	"github.com/jamesrr39/goutil/open"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	envLocation := kingpin.Arg("env-location", "filepath to the environment").Default("env-up-conf.yaml").String()
	port := kingpin.Flag("p", "port").Default("9010").Uint16()
	kingpin.Parse()

	environmentRepo, err := repository.NewEnvironmentRepository(*envLocation)
	if err != nil {
		log.Fatalln(err)
	}

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Compress(gzip.DefaultCompression))
	mainRouter.Route("/api/", func(r chi.Router) {
		r.Mount("/environment", webservices.NewEnvironmentWebService(environmentRepo))
	})
	mainRouter.Mount("/", staticassets.NewStaticAssetsHandler())

	server := httpextra.NewServerWithTimeouts()
	addr := fmt.Sprintf("localhost:%d", *port)
	server.Addr = addr
	server.Handler = mainRouter
	log.Printf("serving on %q\n", addr)
	go func() {
		time.Sleep(time.Millisecond * 500)
		open.OpenURL(fmt.Sprintf("http://%s", addr))
	}()
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
