package webservices

import (
	"env-up-app/backend/repository"
	"env-up-app/backend/types"
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
	router.Post("/{componentName}/start", ws.handleStartComponent)
	router.Post("/{componentName}/stop", ws.handleStopComponent)
	router.Mount("/logs", NewLogsWebsocketService(environmentRepository))

	return ws
}

func (ws *EnvironmentWebService) handleGet(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ws.environmentRepository.Get())
}

func (ws *EnvironmentWebService) handleStartComponent(w http.ResponseWriter, r *http.Request) {
	componentName := chi.URLParam(r, "componentName")
	err := ws.environmentRepository.Start(componentName)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			http.Error(w, err.Error(), 404)
			return
		default:
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func (ws *EnvironmentWebService) handleStopComponent(w http.ResponseWriter, r *http.Request) {
	componentName := chi.URLParam(r, "componentName")
	err := ws.environmentRepository.Stop(componentName)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			http.Error(w, err.Error(), 404)
			return
		default:
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
