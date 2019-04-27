package models

import (
	"github.com/gin-gonic/gin"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/gorilla/websocket"
)

type Request struct {
	Email  string `json:"email"`
	Action string `json:"action"`
	JWT    string `json:"JWT"`
}

type UserInfo struct {
	SessionID         string           `json:"sessionId"`
	JWT               string           `json:"JWT"`
	Intent            string           `json:"command"`
	Parameter         *structpb.Struct `json:"param"`
	Message           string
	Connection        *websocket.Conn
	WebhookConnection *gin.Context
}
