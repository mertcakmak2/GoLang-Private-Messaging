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
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			fmt.Println("********************************")
			fmt.Println("Message: ", string(msg))

			mapMsg := map[string]string{}
			json.Unmarshal(msg, &mapMsg)

			receiver := strings.Split(strings.Split(q.Request.URL.Path, "/channel/")[1], "/ws")[0]
			fmt.Printf("receiver: %v\n", receiver)

			if q.Request.URL.Path == s.Request.URL.Path || mapMsg["receiver"] == receiver {
				fmt.Printf("q.Request.URL.Path: %v\n", q.Request.URL.Path)
				fmt.Printf("s.Request.URL.Path: %v\n", s.Request.URL.Path)
				fmt.Printf("mapMsg[\"receiver\"]: %v\n", mapMsg["receiver"])
				fmt.Printf("mapMsg[\"sender\"]: %v\n", mapMsg["sender"])
				return true
			}
			return false
		})
	})

	r.Run(":5000")
}
