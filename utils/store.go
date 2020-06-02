// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package utils

import (
	"strconv"
	"strings"

	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule/store"
	"github.com/jasonjoo2010/goschedule/store/memory"
	"github.com/jasonjoo2010/goschedule/store/redis"
	"github.com/sirupsen/logrus"
)

func CreateStore(s types.Storage) store.Store {
	switch s.Type {
	case "memory":
		return memory.New()
	case "redis":
		port := 6379
		addr := s.Address
		pos := strings.IndexRune(addr, ':')
		if pos > 0 {
			addr = s.Address[:pos]
			var err error
			port, err = strconv.Atoi(s.Address[pos+1:])
			if err != nil {
				logrus.Error("Wrong address format: ", s.Address)
				return nil
			}
		}
		return redis.New(s.Namespace, addr, port)
	default:
		logrus.Warn("Unknow type of storage: ", s.Type)
		return nil
	}
}
