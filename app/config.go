package app

import (
	"errors"
	"os"

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

func (app *App) GetConfig() *types.DashboardConfig {
	return app.conf
}

func (app *App) UpdateStorageConfig(s types.Storage) error {
	app.conf.Storage = s
	data, err := yaml.Marshal(&app.conf)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(app.confFile, os.O_CREATE|os.O_WRONLY, 0644)
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
		return errors.New("Parse storage configuration failed")
	}
	app.conf = conf
	if app.store != nil {
		app.store.Close()
		app.store = nil
	}
	app.store = utils.CreateStore(app.conf.Storage)
	return nil
}
