package webservices

import (
	"env-up-app/backend/types"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type LogsWebsocketService struct {
	upgrader websocket.Upgrader
	chi.Router
}

func NewLogsWebsocketService() *LogsWebsocketService {
	router := chi.NewRouter()
	service := &LogsWebsocketService{websocket.Upgrader{}, router}
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
	i := 1
	for {
		time.Sleep(time.Second)
		c.WriteJSON(types.NewLogMessage(
			&types.Component{
				Name: "app 1",
			},
			types.PipeStdout,
			fmt.Sprintf("#%d", i),
		))
		i++
	}
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
