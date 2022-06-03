// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package app

import (
	"errors"
	"os"
	"strings"

	"github.com/jasonjoo2010/goschedule-console/controller"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule-console/utils"
	"gopkg.in/yaml.v3"
)

func parseConfig(filepath string) (*types.DashboardConfig, error) {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	data := make([]byte, 0, 128)
	buf := make([]byte, 2)
	for {
		read, err := f.Read(buf)
		if read == 0 || err != nil {
			break
		}
		data = append(data, buf[:read]...)
	}
	conf := types.DashboardConfig{}
	err = yaml.Unmarshal(data, &conf)
	return &conf, nil
}

func (app *App) UpdateStorageConfig(s types.Storage) error {
	app.Conf.Storage = s
	data, err := yaml.Marshal(&app.Conf)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(app.confFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Write(data)
	return app.ReloadConfig()
}

func (app *App) InitConfig(filepath string) error {
	app.confFile = filepath
	return app.ReloadConfig()
}

func (app *App) ReloadConfig() error {
	conf, err := parseConfig(app.confFile)
	if err != nil {
		return err
	}
	if !types.VerifyStorage(conf.Storage) {
		return errors.New("parse storage configuration failed")
	}
	app.Conf = conf
	if app.Store != nil {
		app.Store.Close()
		app.Store = nil
	}
	app.Store = utils.CreateStore(app.Conf.Storage)
	if app.Conf.Base != "" {
		if !strings.HasSuffix(app.Conf.Base, "/") {
			app.Conf.Base += "/"
		}
		controller.SetBasePath(app.Conf.Base)
	}
	return nil
}
