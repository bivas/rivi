package server

import (
	"fmt"
	"strings"

	"github.com/bivas/rivi/bot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/gin-gonic/gin.v1"
)

type BotServer struct {
	Bot  bot.Bot
	Uri  string
	Port int

	engine *gin.Engine
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
	gin.SetMode(gin.ReleaseMode)
	server.engine = gin.Default()
}

func (server *BotServer) registerMetrics() {
	registry := prometheus.NewRegistry()
	registry.Register(prometheus.NewGoCollector())
	server.engine.Any("/metrics", gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))
}

func (server *BotServer) registerDefaultHandler() {
	server.engine.GET("/", func(c *gin.Context) {
		c.String(200, "Running RiviBot")
	})
	if server.Uri != "/" {
		server.engine.GET(server.Uri, func(c *gin.Context) {
			c.String(200, "Running RiviBot")
		})
	}
}

func (server *BotServer) Run() error {
	server.initDefaults()
	server.registerMetrics()
	server.registerDefaultHandler()
	server.engine.POST(server.Uri, func(c *gin.Context) {
		result := server.Bot.HandleEvent(c.Request)
		c.JSON(200, result)
	})
	return server.engine.Run(fmt.Sprintf(":%d", server.Port))
}
