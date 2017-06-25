package server

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"gopkg.in/gin-gonic/gin.v1"
	"strings"
)

type BotServer struct {
	Bot  bot.Bot
	Uri  string
	Port int
}

func (server *BotServer) initDefaults() {
	if server.Port == 0 {
		server.Port = 8080
	}
	if server.Uri == "" {
		server.Uri = "/"
	} else if !strings.HasPrefix(server.Uri, "/") {
		server.Uri = "/" + server.Uri
	}
}

func (server *BotServer) Run() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.String(200, "Running rivi")
	})
	if server.Uri != "/" {
		engine.GET(server.Uri, func(c *gin.Context) {
			c.String(200, "Running rivi")
		})
	}
	engine.POST(server.Uri, func(c *gin.Context) {
		server.Bot.HandleEvent(c.Request)
		c.Status(200)
	})
	if server.Port == 0 {
		server.Port = 8080
	}
	return engine.Run(fmt.Sprintf(":%d", server.Port))
}
