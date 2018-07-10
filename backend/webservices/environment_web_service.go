package webservices

import (
	"env-up-app/backend/repository"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type EnvironmentWebService struct {
	chi.Router
	environmentRepository *repository.EnvironmentRepository
}

func NewEnvironmentWebService(environmentRepository *repository.EnvironmentRepository) *EnvironmentWebService {
	router := chi.NewRouter()

	ws := &EnvironmentWebService{router, environmentRepository}
	router.Get("/", ws.handleGet)

	return ws
}

func (ws *EnvironmentWebService) handleGet(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ws.environmentRepository.Get())
}
