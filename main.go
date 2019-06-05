package main

import (
	endpoints "github.com/flow_server/Endpoints"
	models "github.com/flow_server/Models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*Application entry point, creates Websocket endpoint & Webhook endpoint*/
func main() {
	var err error
	r := gin.Default()

	connections := make(map[string]*models.UserInfo)

	//Dialogflow endpoint
	r.POST("/3000", func(c *gin.Context) {
		endpoints.HandleWebhook(c, &connections)
	})

	//Client endpoint
	r.GET("/4000", func(c *gin.Context) {
		endpoints.WebsocketHandler(c.Writer, c.Request, &connections)
	})

	if err = r.Run(":9090"); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
