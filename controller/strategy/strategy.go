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
)

func Init(engine *gin.Engine) {
	group := engine.Group("/strategy")

	group.GET("/index", indexHandler)
	group.POST("/create", createHandler)
	group.GET("/pause", pauseHandler)
	group.GET("/resume", resumeHandler)
	group.POST("/save", saveHandler)
}

func indexHandler(c *gin.Context) {
	list, _ := app.Instance().Store.GetStrategies()
	c.HTML(http.StatusOK, "strategy/index.html", controller.DataWithSession(gin.H{
		"strategies": list,
	}))
}

func createHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	id := strings.TrimSpace(c.PostForm("id"))
	kind := c.PostForm("kind")
	bind := strings.TrimSpace(c.PostForm("bind"))
	limit, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("limit")))
	total, _ := strconv.Atoi(strings.TrimSpace(c.PostForm("total")))
	parameter := strings.TrimSpace(c.PostForm("parameter"))
	target := strings.TrimSpace(c.PostForm("target"))
	store := app.Instance().Store
	if id == "" {
		resp.Err(1, "ID cannot be empry")
		return
	}
	if kind == "" {
		resp.Err(1, "Kind is incorrect")
		return
	}
	if bind == "" {
		resp.Err(1, "Bind cannot be empry")
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
	s, err := store.GetStrategy(id)
	if err == nil && s != nil {
		resp.Err(2, "Strategy in same ID has already existed")
		return
	}
	if s != nil {
		resp.Err(3, "Strategy's ID should be unique")
		return
	}
	if target == "" {
		target = "127.0.0.1"
	}
	strategyKind := utils.ToStrategyKind(kind)
	if strategyKind == definition.UnknowKind {
		resp.Err(4, "Kind is illegal")
		return
	}
	store.CreateStrategy(&definition.Strategy{
		Id:                   id,
		MaxOnSingleScheduler: limit,
		Total:                total,
		Kind:                 strategyKind,
		Bind:                 bind,
		Parameter:            parameter,
		IpList:               strings.Split(target, ","),
		Enabled:              true,
	})
}

func saveHandler(c *gin.Context) {
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
