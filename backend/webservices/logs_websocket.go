package webservices

import (
	"env-up-app/backend/repository"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type LogsWebsocketService struct {
	environmentRepository *repository.EnvironmentRepository
	upgrader              websocket.Upgrader
	chi.Router
}

func NewLogsWebsocketService(environmentRepository *repository.EnvironmentRepository) *LogsWebsocketService {
	router := chi.NewRouter()
	service := &LogsWebsocketService{environmentRepository, websocket.Upgrader{}, router}
	router.HandleFunc("/", service.handleWebsocket)
	return service
}

func (s *LogsWebsocketService) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	listener := s.environmentRepository.GetLogMessageChanListener()
	defer listener.Close()

	closeChan := make(chan struct{})

	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
					println("close no status received")
					closeChan <- struct{}{}
					return
				}
				log.Printf("error reading message: %q\n", err)
				return
			}
			log.Printf("got message: type: %d. Message: %q\n", mt, message)
		}
	}()

	for {
		select {
		case message := <-listener.Chan:
			err := c.WriteJSON(message)
			if err != nil {
				log.Printf("failed to marshal to json and write to websocket. Error: %q. Object: %v\n", err, message)
				continue
			}
		case <-closeChan:
			return
		}
	}
}
