package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	security "github.com/flow_server/Security"
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

			valid := security.ValidateClientMessage(request.JWT)

			if !valid {
				return
			}

			_, bodyStructValidationError := security.ValidateJwt(request.JWT)
			utils.ErrorInvalidCheck(bodyStructValidationError)
			if bodyStructValidationError != nil {
				return
			}

			if (len(request.Email) > 0) && (request.Action == "Sign In") {
				fmt.Println("Channel Created!")
				email := strings.ToLower(request.Email)
				connec := *connections

				if connec[email] == nil {
					fmt.Println("New Connection Created! Websocket")
					connec[email] = &models.UserInfo{
						Connection: conn,
					}
				} else {
					fmt.Println("Connection found, adding conn Websocket")
					payload := connec[email]
					payload.Connection = conn

				}
			}
		}

	}(conn)
}
