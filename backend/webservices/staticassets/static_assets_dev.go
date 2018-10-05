// +build !prod

package staticassets

import "net/http"

func NewStaticAssetsHandler() http.Handler {
	return http.FileServer(http.Dir("frontend-web"))
}
