//+build prod
package webservices

import (
	"log"
	"net/http"

	_ "env-up-app/build/client/statik"

	"github.com/rakyll/statik/fs"
)

func NewStaticAssetsHandler() http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(statikFS)
}
