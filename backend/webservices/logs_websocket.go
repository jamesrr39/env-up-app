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

	for {
		message := <-s.environmentRepository.LogMessageChan
		err := c.WriteJSON(message)
		if err != nil {
			log.Printf("failed to marshal to json and write to websocket. Error: %q. Object: %q\n", err, message)
			continue
		}
	}
	// i := 1
	// for {
	// 	time.Sleep(time.Second)
	// 	c.WriteJSON(types.NewLogMessage(
	// 		&types.Component{
	// 			Name: "app 1",
	// 		},
	// 		types.PipeStdout,
	// 		fmt.Sprintf("#%d", i),
	// 	))
	// 	i++
	// }
	// for {
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	err = c.WriteMessage(mt, message)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
}
