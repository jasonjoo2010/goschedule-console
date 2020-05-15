package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule-console/utils"
)

func Init(engine *gin.Engine) {
	configGroup := engine.Group("/config")

	configGroup.GET("/modify", modifyHandler)
	configGroup.POST("/save", saveHandler)
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
