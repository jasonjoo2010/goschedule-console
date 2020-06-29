// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package strategy

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule-console/utils"
	"github.com/jasonjoo2010/goschedule/core/definition"
	storepkg "github.com/jasonjoo2010/goschedule/store"
	"github.com/robfig/cron"
)

func Init(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/strategy/index")
	})
	group := engine.Group("/strategy")

	group.GET("/index", indexHandler)
	group.GET("/get", getHandler)
	group.POST("/create", createHandler)
	group.GET("/pause", pauseHandler)
	group.GET("/resume", resumeHandler)
	group.POST("/save", saveHandler)
	group.GET("/remove", removeHandler)
	group.GET("/info", infoHandler)
}

func indexHandler(c *gin.Context) {
	list, _ := app.Instance().Store.GetStrategies()
	c.HTML(http.StatusOK, "strategy/index.html", controller.DataWithSession(gin.H{
		"strategies": list,
	}))
}

func getHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No strategy specified")
		return
	}
	s := app.Instance().Store
	strategy, err := s.GetStrategy(id)
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	resp["strategy"] = strategy
}

func infoHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No strategy specified")
		return
	}
	s := app.Instance().Store
	runtimes, err := s.GetStrategyRuntimes(id)
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	resp["runtimes"] = runtimes
}

func removeHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No strategy specified")
		return
	}
	s := app.Instance().Store
	_, err := s.GetStrategy(id)
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	s.RemoveStrategy(id)
}

func checkStrategy(resp types.JsonResponse, c *gin.Context) (strategy *definition.Strategy, result bool) {
	id := strings.TrimSpace(c.PostForm("id"))
	kind := c.PostForm("kind")
	bind := strings.TrimSpace(c.PostForm("bind"))
	limit, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("limit")))
	total, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("total")))
	parameter := strings.TrimSpace(c.PostForm("parameter"))
	target := strings.TrimSpace(c.PostForm("target"))
	cronBegin := strings.TrimSpace(c.PostForm("cronBegin"))
	cronEnd := strings.TrimSpace(c.PostForm("cronEnd"))
	if id == "" {
		resp.Err(1, "ID cannot be empty")
		return
	}
	if kind == "" {
		resp.Err(1, "Kind is incorrect")
		return
	}
	if bind == "" {
		resp.Err(1, "Bind cannot be empty")
		return
	}
	if total < 1 {
		resp.Err(1, "Total should be 1 at least")
		return
	}
	if limit < 0 || limit > total {
		resp.Err(1, "Limit should be [0, total]")
		return
	}
	if target == "" {
		target = "127.0.0.1"
	}
	cronParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if cronBegin != "" {
		_, err := cronParser.Parse(cronBegin)
		if err != nil {
			resp.Err(1, "CronBegin is incorrect.")
			return
		}
	}
	if cronEnd != "" {
		_, err := cronParser.Parse(cronEnd)
		if err != nil {
			resp.Err(1, "CronEnd is incorrect.")
			return
		}
	}
	strategyKind := utils.ToStrategyKind(kind)
	if strategyKind == definition.UnknowKind {
		resp.Err(4, "Kind is illegal")
		return
	}
	strategy = &definition.Strategy{
		Id:                   id,
		MaxOnSingleScheduler: limit,
		Total:                total,
		Kind:                 strategyKind,
		Bind:                 bind,
		Parameter:            parameter,
		Enabled:              false,
		CronBegin:            cronBegin,
		CronEnd:              cronEnd,
	}
	targets := strings.Split(target, ",")
	strategy.IpList = make([]string, 0, len(targets))
	for _, t := range targets {
		if len(t) < 1 {
			continue
		}
		strategy.IpList = append(strategy.IpList, t)
	}
	result = true
	return
}

func createHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	store := app.Instance().Store
	if strategy, ok := checkStrategy(resp, c); ok {
		s, err := store.GetStrategy(strategy.Id)
		if err != nil && err != storepkg.NotExist {
			resp.Err(2, "Fail to retrieve data from store")
			return
		}
		if s != nil {
			resp.Err(2, "Strategy in same ID has already existed")
			return
		}
		store.CreateStrategy(strategy)
	}
}

func saveHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	store := app.Instance().Store
	if strategy, ok := checkStrategy(resp, c); ok {
		s, err := store.GetStrategy(strategy.Id)
		if err != nil && err != storepkg.NotExist {
			resp.Err(2, "Fail to retrieve data from store")
			return
		}
		if s == nil {
			resp.Err(2, "Strategy doesn't exist")
			return
		}
		strategy.Enabled = s.Enabled
		store.UpdateStrategy(strategy)
	}
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
	strategy, err := store.GetStrategy(id)
	if err != nil || strategy == nil {
		resp.Err(2, "Strategy cannot be found")
		return
	}
	strategy.Enabled = enabled
	store.UpdateStrategy(strategy)
}

func pauseHandler(c *gin.Context) {
	controlHandler(false, c)
}

func resumeHandler(c *gin.Context) {
	controlHandler(true, c)
}
