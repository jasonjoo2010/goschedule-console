// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package task

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule-console/utils"
	"github.com/jasonjoo2010/goschedule/definition"
	storepkg "github.com/jasonjoo2010/goschedule/store"
)

type info struct {
	ConfigVersion int64
	Runtimes      []*definition.TaskRuntime
	Assignments   []*definition.TaskAssignment
}

func Init(engine *gin.RouterGroup) {
	group := engine.Group("/task")

	group.GET("/index", indexHandler)
	group.GET("/get", getHandler)
	group.GET("/info", infoHandler)
	group.POST("/create", createHandler)
	group.POST("/save", saveHandler)
	group.GET("/remove", removeHandler)
}

func indexHandler(c *gin.Context) {
	list, _ := app.Instance().Store.GetTasks()
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	c.HTML(http.StatusOK, "task/index.html", controller.DataWithSession(gin.H{
		"tasks": list,
	}))
}

func getHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No task specified")
		return
	}
	s := app.Instance().Store
	task, err := s.GetTask(id)
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	resp["task"] = task
}

func infoHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No task specified")
		return
	}
	s := app.Instance().Store
	task, err := s.GetTask(id)
	if err != nil {
		resp.Err(1, "Specific task could not be found")
		return
	}
	strategies, _ := s.GetStrategies()
	infoMap := make(map[string]info)
	strategyMap := make(map[string]*definition.Strategy)
	for _, strategy := range strategies {
		if strategy.Kind != definition.TaskKind || strategy.Bind != task.ID {
			continue
		}
		strategyMap[strategy.ID] = strategy
		runtimes, _ := s.GetTaskRuntimes(strategy.ID, task.ID)
		assignments, _ := s.GetTaskAssignments(strategy.ID, task.ID)
		version, _ := s.GetTaskItemsConfigVersion(strategy.ID, task.ID)
		infoMap[strategy.ID] = info{
			Runtimes:      runtimes,
			Assignments:   assignments,
			ConfigVersion: version,
		}
	}
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	resp["strategies"] = strategyMap
	resp["info"] = infoMap
}

func removeHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := c.Query("id")
	if id == "" {
		resp.Err(1, "No task specified")
		return
	}
	s := app.Instance().Store
	_, err := s.GetTask(id)
	if err != nil {
		resp.Err(2, err.Error())
		return
	}
	s.RemoveTask(id)
}

func checkTask(resp types.JsonResponse, c *gin.Context) (task *definition.Task, result bool) {
	id := strings.TrimSpace(c.PostForm("id"))
	model := c.PostForm("model")
	bind := strings.TrimSpace(c.PostForm("bind"))
	fetch, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("fetch")))
	batch, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("batch")))
	executor, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("executor")))
	maxtaskitem, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("maxtaskitem")))
	parameter := strings.TrimSpace(c.PostForm("parameter"))
	taskitem := strings.TrimSpace(c.PostForm("taskitem"))
	interval, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("interval")))
	intervalNoData, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("intervalNoData")))
	heartbeat, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("heartbeat")))
	death, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("death")))
	if id == "" {
		resp.Err(1, "ID cannot be empty")
		return
	}
	if model == "" {
		resp.Err(1, "Model is incorrect")
		return
	}
	if bind == "" {
		resp.Err(1, "Bind cannot be empty")
		return
	}
	if fetch < 1 {
		resp.Err(1, "Fetch should be 1 at least")
		return
	}
	if batch < 0 || batch > fetch {
		resp.Err(1, "Batch should be [1, FetchCount]")
		return
	}
	if executor < 1 {
		resp.Err(1, "Executor count should be 1 at least")
		return
	}
	if maxtaskitem < 0 {
		resp.Err(1, "Maximum task item per worker should be nonnegative")
		return
	}
	if heartbeat < 1000 {
		resp.Err(1, "Heartbeat is overfrequency it maybe inefficiency")
		return
	}
	if death < 2*heartbeat {
		resp.Err(1, "Death timeout should be at least twice the interval of heartbeat")
		return
	}
	items := utils.ParseTaskItems(taskitem)
	if len(items) < 1 {
		resp.Err(1, "No valid task item found")
		return
	}
	if interval < 0 {
		interval = 0
	}
	if intervalNoData < 0 {
		intervalNoData = 0
	}
	taskModel := definition.Normal
	if model == "stream" {
		taskModel = definition.Stream
	}
	task = &definition.Task{
		ID:                id,
		Model:             taskModel,
		Bind:              bind,
		Parameter:         parameter,
		FetchCount:        fetch,
		BatchCount:        batch,
		ExecutorCount:     executor,
		MaxTaskItems:      maxtaskitem,
		Items:             items,
		Interval:          interval,
		IntervalNoData:    intervalNoData,
		HeartbeatInterval: heartbeat,
		DeathTimeout:      death,
	}
	result = true
	return
}

func createHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	store := app.Instance().Store
	if task, ok := checkTask(resp, c); ok {
		s, err := store.GetTask(task.ID)
		if err != nil && err != storepkg.NotExist {
			resp.Err(2, "Fail to retrieve data from store: "+err.Error())
			return
		}
		if s != nil {
			resp.Err(2, "Task in same ID has already existed")
			return
		}
		store.CreateTask(task)
	}
}

func saveHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	store := app.Instance().Store
	if task, ok := checkTask(resp, c); ok {
		s, err := store.GetTask(task.ID)
		if err != nil && err != storepkg.NotExist {
			resp.Err(2, "Fail to retrieve data from store: "+err.Error())
			return
		}
		if s == nil {
			resp.Err(2, "Task doesn't exist")
			return
		}
		store.UpdateTask(task)
	}
}
