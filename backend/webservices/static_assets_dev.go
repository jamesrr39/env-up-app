//+build !prod

package webservices

import "net/http"

func NewStaticAssetsHandler() http.Handler {
	return http.FileServer(http.Dir("frontend-web"))
}
