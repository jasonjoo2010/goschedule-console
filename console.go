// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jasonjoo2010/goschedule-console/app"
)

var (
	configFile string
	listenPort int
	debug      bool
)

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "Configuration file")
	flag.IntVar(&listenPort, "p", 8000, "Serving port")
	flag.BoolVar(&debug, "v", false, "Debug mode")
}

func main() {
	flag.Parse()
	err := app.Instance().InitConfig(configFile)
	if err != nil {
		return
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := InitEngine()
	engine.Run(":" + strconv.Itoa(listenPort))
}
