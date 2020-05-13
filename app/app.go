package app

import (
	"sync"

	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule/store"
)

var (
	lock     = sync.Mutex{}
	instance *App
)

type App struct {
	confFile string
	conf     *types.DashboardConfig
	store    store.Store
}

func Instance() *App {
	// likely
	if instance != nil {
		return instance
	}
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = &App{}
	}
	return instance
}
