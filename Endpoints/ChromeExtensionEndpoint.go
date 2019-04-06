package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/flow_server/Utils"

	models "github.com/flow_server/Models"
	"github.com/gorilla/websocket"
)

var connectionUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*WebsocketHandler handles all client side message logic, parsing
all client messages and establish user profile in server.*/
func WebsocketHandler(w http.ResponseWriter, r *http.Request, connections *map[string]*models.UserInfo) {
	var conn, err = connectionUpgrader.Upgrade(w, r, nil)
	utils.ErrorNilCheck(err)

	conn.WriteJSON(models.UserInfo{
		Message: "Connection successfully established"})

	clientToServerMessageHandler(connections, conn)

}

func clientToServerMessageHandler(connections *map[string]*models.UserInfo, conn *websocket.Conn) {

	go func(conn *websocket.Conn) {
		_, p, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
		} else {
			request := models.Request{}

			error := json.Unmarshal(p, &request)
			utils.ErrorNilCheck(error)

			if (len(request.Email) > 0) && (request.Action == "Sign In") {
				fmt.Println("Channel Created!")
				connec := *connections

				if connec[request.Email] == nil {
					fmt.Println("New Connection Created! Websocket")
					connec[request.Email] = &models.UserInfo{
						Connection: conn,
					}
				} else {
					fmt.Println("Connection found, adding conn Websocket")
					payload := connec[request.Email]
					payload.Connection = conn

				}
			}
		}

	}(conn)
}
