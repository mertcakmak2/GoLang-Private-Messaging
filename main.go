package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/channel/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "chan.html")
	})

	r.GET("/channel/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		go sendMessageToUser(m, s, msg)
		fmt.Println("handle message")
	})

	r.Run(":5000")
}

func sendMessageToUser(m *melody.Melody, s *melody.Session, msg []byte) {
	m.BroadcastFilter(msg, func(q *melody.Session) bool {
		message := convertByteToMessage(msg)
		receiver := getReceiver(q.Request.URL.Path)
		if q.Request.URL.Path == s.Request.URL.Path || message["receiver"] == receiver {
			return true
		}
		return false
	})
}

func convertByteToMessage(msg []byte) map[string]string {
	message := map[string]string{}
	json.Unmarshal(msg, &message)
	return message
}

func getReceiver(path string) string {
	receiver := strings.Split(strings.Split(path, "/channel/")[1], "/ws")[0]
	return receiver
}
