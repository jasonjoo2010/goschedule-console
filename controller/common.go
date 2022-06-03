// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package controller

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule/log"
)

var basePath atomic.Value

func init() {
	basePath.Store("/")
}

func SetBasePath(path string) {
	log.Infof("Set base path to %s", path)
	basePath.Store(path)
}

func DataWithSession(data map[string]interface{}) map[string]interface{} {
	if gin.Mode() == gin.ReleaseMode {
		data["ENV"] = "production"
	} else {
		data["ENV"] = "development"
	}
	data["basePath"] = basePath.Load()
	return data
}
