package endpoints

import (
	"fmt"
	"net/http"

	utils "github.com/flow_server/Utils"

	models "github.com/flow_server/Models"
	security "github.com/flow_server/Security"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

/*HandleWebhook handles all request from Dialogflow. Obtain message,
parse JWT tokens and parse commands into Struct objects. Passes result
to wsHandler.*/
func HandleWebhook(c *gin.Context, connections *map[string]*models.UserInfo) {
	var err error
	var unmarshall jsonpb.Unmarshaler
	unmarshall.AllowUnknownFields = true

	wr := dialogflow.WebhookRequest{}

	if err = unmarshall.Unmarshal(c.Request.Body, &wr); err != nil {
		logrus.WithError(err).Error("Couldn't Unmarshal request to jsonpb")
		c.Status(http.StatusBadRequest)
		return
	}
	//fmt.Print(wr.String())
	queryParameters := wr.GetQueryResult().GetParameters()
	queryIntent := wr.GetQueryResult().GetIntent().GetDisplayName()
	tokenPayloadResponse := wr.GetOriginalDetectIntentRequest().GetPayload().GetFields()
	userStruct := tokenPayloadResponse["user"]
	idTokenString := userStruct.GetStructValue().GetFields()["idToken"].GetStringValue()

	bodyStruct, bodyStructValidationError := security.ValidateJwt(idTokenString)

	utils.ErrorInvalidCheck(bodyStructValidationError)

	if bodyStructValidationError != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	con := *connections
	// Note: If client connection is persistent, we need to throw an error if bodyStruct.email
	// is not found, meaning there is no client to reach out to. We should respond to dialogflow telling
	// user that connection is not establish, check client.
	// If client terminates or shuts down, we need to remove the registered information

	if con[bodyStruct.Email] == nil {
		fmt.Println("Client connection is not established")
		fullfillment := dialogflow.WebhookResponse{
			FulfillmentText: "Client connection is not established, please check your chrome extension to see if it is turned on.",
		}
		c.JSON(http.StatusOK, fullfillment)
		return
	}

	fmt.Println("Connection found writing to JSON, Webhook")
	payload := *con[bodyStruct.Email]

	if payload.SessionID != wr.Session {
		payload.SessionID = wr.Session
		payload.Intent = ""
		payload.Parameter = nil
		payload.Connection.WriteJSON(payload)
	}

	payload.Intent = queryIntent
	payload.Parameter = queryParameters

	payload.JWT = security.CreateJwtToken()
	payload.Connection.WriteJSON(payload)

	fullfillment := dialogflow.WebhookResponse{
		FulfillmentText: " ",
	}

	c.JSON(http.StatusOK, fullfillment)
}

// func terminateAgentConnection(userEmail string, connections *map[string]*models.UserInfo, c *gin.Context) {
// 	c.JSON(http.StatusOK)
// }
