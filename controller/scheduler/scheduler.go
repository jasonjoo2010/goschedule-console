// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package scheduler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
)

func Init(engine *gin.RouterGroup) {
	group := engine.Group("/scheduler")

	group.GET("/index", indexHandler)
	group.GET("/stop", stopHandler)
	group.GET("/start", startHandler)
}

func indexHandler(c *gin.Context) {
	list, _ := app.Instance().Store.GetSchedulers()
	c.HTML(http.StatusOK, "scheduler/index.html", controller.DataWithSession(gin.H{
		"schedulers": list,
	}))
}

func controlHandler(enabled bool, c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	store := app.Instance().Store
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "id cannot be empty")
		return
	}
	scheduler, err := store.GetScheduler(id)
	if err != nil || scheduler == nil {
		resp.Err(2, "Scheduler cannot be found")
		return
	}
	scheduler.Enabled = enabled
	store.RegisterScheduler(scheduler)
}

func stopHandler(c *gin.Context) {
	controlHandler(false, c)
}

func startHandler(c *gin.Context) {
	controlHandler(true, c)
}
