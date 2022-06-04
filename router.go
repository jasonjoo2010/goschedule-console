// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/controller/config"
	"github.com/jasonjoo2010/goschedule-console/controller/scheduler"
	"github.com/jasonjoo2010/goschedule-console/controller/strategy"
	"github.com/jasonjoo2010/goschedule-console/controller/task"
)

func add(numbers ...interface{}) interface{} {
	if len(numbers) < 1 {
		return 0
	}
	if len(numbers) < 2 {
		return numbers[0]
	}
	var total int64
	var totalf float64
	float_result := false
	for _, num := range numbers {
		switch n := num.(type) {
		case int8:
			total += int64(n)
			totalf += float64(n)
		case int:
			total += int64(n)
			totalf += float64(n)
		case int16:
			total += int64(n)
			totalf += float64(n)
		case int32:
			total += int64(n)
			totalf += float64(n)
		case int64:
			total += n
			totalf += float64(n)
		case uint:
			total += int64(n)
			totalf += float64(n)
		case uint8:
			total += int64(n)
			totalf += float64(n)
		case uint16:
			total += int64(n)
			totalf += float64(n)
		case uint32:
			total += int64(n)
			totalf += float64(n)
		case uint64:
			total += int64(n)
			totalf += float64(n)
		case float32:
			float_result = true
			totalf += float64(n)
		case float64:
			float_result = true
			totalf += n
		}
	}
	if float_result {
		return totalf
	}
	return total
}

func timestampMillis(millis int64) string {
	tm := time.Unix(millis/1000, (millis%1000)*time.Hour.Milliseconds())
	return tm.Format("2006-01-02 15:04:05")
}

func timestamp(ts int64) string {
	tm := time.Unix(ts, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func storageChecker(c *gin.Context) {
	switch c.FullPath() {
	case controller.GetBasePath() + "config/modify", controller.GetBasePath() + "config/save":
		return
	default:
	}
	store := app.Instance().Store
	if store == nil {
		c.Redirect(302, app.Instance().Conf.Base+"config/modify")
		c.Abort()
		return
	}
}

func InitEngine() *gin.Engine {
	engine := gin.New()
	engine.FuncMap["add"] = add
	engine.FuncMap["timestamp"] = timestamp
	engine.FuncMap["timestampMillis"] = timestampMillis
	engine.LoadHTMLGlob("templates/**/*")
	rootGroup := engine.Group(controller.GetBasePath())
	rootGroup.Static("/css", "static/css")
	rootGroup.Static("/js", "static/js")
	rootGroup.Use(gin.Recovery())
	rootGroup.Use(storageChecker)
	rootGroup.Use(func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.Set("ENV", "production")
		} else {
			c.Set("ENV", "development")
		}
	})
	config.Init(rootGroup)
	strategy.Init(rootGroup)
	scheduler.Init(rootGroup)
	task.Init(rootGroup)
	return engine
}
