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
	router.Get("/{filePath}", ws.handleGet)

	return ws
}

func (ws *EnvironmentWebService) handleGet(w http.ResponseWriter, r *http.Request) {
	filePath := chi.URLParam(r, "filePath")
	environment, err := ws.environmentRepository.Get(filePath)
	if err != nil {
		http.Error(w, err.Error(), 500) // TODO better error code
		return
	}

	render.JSON(w, r, environment)
}
