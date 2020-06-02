// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package types

import (
	"github.com/sirupsen/logrus"
)

var STORAGE_TYPES [4]string = [...]string{"memory", "redis", "zookeeper", "jdbc"}

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
	logrus.Error("Type of storage is unknown: ", s.Type)
	return false
}
