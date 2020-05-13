package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/controller/config"
)

func InitEngine() *gin.Engine {
	engine := gin.New()
	engine.Static("/css", "static/css")
	engine.Static("/js", "static/js")
	engine.LoadHTMLGlob("templates/**/*")
	engine.Use(gin.Recovery())
	engine.Use(func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.Set("ENV", "production")
		} else {
			c.Set("ENV", "development")
		}
	})
	config.Init(engine)
	return engine
}
