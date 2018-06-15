package webservices

import (
	"encoding/json"
	"env-up-app/backend/repository"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ConfigWebService struct {
	chi.Router
	configRepository *repository.ConfigRepository
}

func NewConfigWebService(configRepository *repository.ConfigRepository) *ConfigWebService {
	router := chi.NewRouter()

	ws := &ConfigWebService{router, configRepository}
	router.Get("/", ws.handleGet)
	router.Post("/", ws.handleCreateEnvironment)

	return ws
}

func (ws *ConfigWebService) handleGet(w http.ResponseWriter, r *http.Request) {
	config := ws.configRepository.Get()

	render.JSON(w, r, config)
}

type postEnvironmentRequestBody struct {
	EnvironmentFilePath string `json:"environmentFilePath"`
}

func (ws *ConfigWebService) handleCreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var requestBody postEnvironmentRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = ws.configRepository.AddEnvironmentPath(requestBody.EnvironmentFilePath)
	if err != nil {
		http.Error(w, err.Error(), 500) // TODO better error
		return
	}

	w.WriteHeader(204)
}
