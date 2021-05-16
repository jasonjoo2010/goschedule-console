// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package types

import (
	"github.com/jasonjoo2010/goschedule/log"
)

var STORAGE_TYPES [6]string = [...]string{"memory", "redis", "zookeeper", "database", "etcdv2", "etcdv3"}

// DashboardConfig respresents structure of config.xml
type DashboardConfig struct {
	Storage Storage
}

type Storage struct {
	Type      string
	Address   string
	Username  string
	Password  string
	Namespace string
}

func VerifyStorage(s Storage) bool {
	if s.Type == "" {
		return true
	}
	for _, t := range STORAGE_TYPES {
		if t == s.Type {
			return true
		}
	}
	log.Errorf("Type of storage is unknown: %s", s.Type)
	return false
}
