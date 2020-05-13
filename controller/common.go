package controller

import "github.com/gin-gonic/gin"

func DataWithSession(data map[string]interface{}) map[string]interface{} {
	if gin.Mode() == gin.ReleaseMode {
		data["ENV"] = "production"
	} else {
		data["ENV"] = "development"
	}
	return data
}
