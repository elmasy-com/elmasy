package main

import (
	"fmt"
	"time"

	"github.com/elmasy-com/elmasy/config"
	"github.com/elmasy-com/elmasy/router/api/ip"
	"github.com/elmasy-com/elmasy/router/api/protocol/dns"
	"github.com/elmasy-com/elmasy/router/api/protocol/tls"
	randomip "github.com/elmasy-com/elmasy/router/api/random/ip"
	randomport "github.com/elmasy-com/elmasy/router/api/random/port"
	"github.com/elmasy-com/elmasy/router/api/scan/port"

	"github.com/gin-gonic/gin"

	ginstatic "github.com/gin-contrib/static"
)

func logFormat(param gin.LogFormatterParams) string {

	return fmt.Sprintf("%s - [%s] \"%s %s\" %s %d %s \"%s\"\n%s",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func SetupRouter() *gin.Engine {

	gin.DisableConsoleColor()

	if !config.GlobalConfig.Verbose {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logFormat))

	engine.SetTrustedProxies(config.GlobalConfig.TrustedProxies)

	//engine.StaticFS("/doc", doc.MustFS())
	engine.Use(ginstatic.ServeRoot("/", "./static"))

	api := engine.Group("/api")
	{
		api.GET("/ip", ip.Get)
		api.GET("/random/ip/:version", randomip.Get)
		api.GET("/random/port", randomport.Get)
		api.GET("/protocol/dns/:type/:name", dns.Get)
		api.GET("/protocol/tls", tls.Get)
		api.GET("/scan/port", port.Get)
	}

	return engine
}
