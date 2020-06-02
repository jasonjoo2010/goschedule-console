// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

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
