package config

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule-console/utils"
	"github.com/jasonjoo2010/goschedule/core/definition"
)

const (
	PREFIX_STRATEGY = "Strategy:"
	PREFIX_TASK     = "Task:"
)

func Init(engine *gin.Engine) {
	configGroup := engine.Group("/config")

	configGroup.GET("/modify", modifyHandler)
	configGroup.POST("/save", saveHandler)
	configGroup.GET("/dump", dumpHandler)
	configGroup.GET("/export", exportHandler)
	configGroup.GET("/import", importHandler)
	configGroup.POST("/importSave", importSaveHandler)
}

func saveHandler(c *gin.Context) {
	resp := types.NewEmptyResponse()
	defer c.JSON(200, resp)
	isTest := c.PostForm("test") == "1"
	storageType := c.PostForm("type")
	addr := c.PostForm("address")
	username := c.PostForm("username")
	password := c.PostForm("password")
	namespace := c.PostForm("namespace")
	if addr == "" && storageType != "memory" {
		resp.Err(2, "Address can't be empty")
		return
	}
	storageConfig := types.Storage{
		Type:      storageType,
		Address:   addr,
		Username:  username,
		Password:  password,
		Namespace: namespace,
	}
	if isTest {
		resp["test"] = true
		s := utils.CreateStore(storageConfig)
		list, err := s.GetSchedulers()
		if err != nil {
			resp.Err(1, "Failed: "+err.Error())
			return
		}
		if len(list) < 1 {
			resp.Err(1, "No schedulers online on specific namespace, please check")
			return
		}
		return
	}
	// save
	application := app.Instance()
	err := application.UpdateStorageConfig(storageConfig)
	if err != nil {
		resp.Err(1, "Update failed: "+err.Error())
	}
}

func modifyHandler(c *gin.Context) {
	conf := app.Instance().Conf.Storage
	c.HTML(http.StatusOK, "config/modify.html", controller.DataWithSession(gin.H{
		"addr":      conf.Address,
		"type":      conf.Type,
		"username":  conf.Username,
		"password":  conf.Password,
		"namespace": conf.Namespace,
	}))
}

func dumpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "config/dump.html", controller.DataWithSession(gin.H{
		"content": app.Instance().Store.Dump(),
	}))
}

func exportHandler(c *gin.Context) {
	store := app.Instance().Store
	b := &strings.Builder{}
	strategies, _ := store.GetStrategies()
	for _, s := range strategies {
		b.WriteString("Strategy:")
		data, _ := json.Marshal(s)
		b.Write(data)
		b.WriteString("\n")
	}
	tasks, _ := store.GetTasks()
	for _, s := range tasks {
		b.WriteString("Task:")
		data, _ := json.Marshal(s)
		b.Write(data)
		b.WriteString("\n")
	}
	c.HTML(http.StatusOK, "config/export.html", controller.DataWithSession(gin.H{
		"exported": b.String(),
	}))
}

func importHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "config/import.html", controller.DataWithSession(gin.H{}))
}

func parseImportContent(content string) ([]*definition.Strategy, []*definition.Task) {
	strategies := make([]*definition.Strategy, 0, 2)
	tasks := make([]*definition.Task, 0, 2)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.Index(line, PREFIX_STRATEGY) == 0 {
			// strategy
			jsonStr := line[len(PREFIX_STRATEGY):]
			strategy := definition.Strategy{}
			err := json.Unmarshal([]byte(jsonStr), &strategy)
			if err != nil {
				continue
			}
			strategy.Enabled = false
			strategies = append(strategies, &strategy)
		} else if strings.Index(line, PREFIX_TASK) == 0 {
			// task
			jsonStr := line[len(PREFIX_TASK):]
			task := definition.Task{}
			err := json.Unmarshal([]byte(jsonStr), &task)
			if err != nil {
				continue
			}
			tasks = append(tasks, &task)
		}
	}
	return strategies, tasks
}

func importSaveHandler(c *gin.Context) {
	msg := types.NewEmptyResponse()
	defer c.JSON(200, msg)
	content := c.PostForm("content")
	if content == "" {
		msg.Err(1, "Empty importing content")
		return
	}
	strategies, tasks := parseImportContent(content)
	strategiesTotal := len(strategies)
	strategiesSuccess := 0
	tasksTotal := len(tasks)
	tasksSuccess := 0
	store := app.Instance().Store
	for _, strategy := range strategies {
		s, _ := store.GetStrategy(strategy.Id)
		if s != nil {
			continue
		}
		err := store.CreateStrategy(strategy)
		if err != nil {
			continue
		}
		strategiesSuccess++
	}
	for _, task := range tasks {
		s, _ := store.GetTask(task.Id)
		if s != nil {
			continue
		}
		err := store.CreateTask(task)
		if err != nil {
			continue
		}
		tasksSuccess++
	}
	msg["strategiesTotal"] = strategiesTotal
	msg["strategiesSuccess"] = strategiesSuccess
	msg["tasksTotal"] = tasksTotal
	msg["tasksSuccess"] = tasksSuccess
}
