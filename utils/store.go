// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package utils

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonjoo2010/goschedule-console/types"
	"github.com/jasonjoo2010/goschedule/log"
	"github.com/jasonjoo2010/goschedule/store"
	"github.com/jasonjoo2010/goschedule/store/database"
	"github.com/jasonjoo2010/goschedule/store/etcdv2"
	"github.com/jasonjoo2010/goschedule/store/etcdv3"
	"github.com/jasonjoo2010/goschedule/store/memory"
	"github.com/jasonjoo2010/goschedule/store/redis"
	"github.com/jasonjoo2010/goschedule/store/zookeeper"
)

func parseAddr(address string, default_port int) (string, int, error) {
	port := default_port
	addr := address
	pos := strings.IndexRune(addr, ':')
	if pos > 0 {
		addr = address[:pos]
		var err error
		port, err = strconv.Atoi(address[pos+1:])
		if err != nil {
			log.Errorf("Wrong address format: %s", address)
			return "", 0, err
		}
	}
	return addr, port, nil
}

func CreateStore(s types.Storage) store.Store {
	switch s.Type {
	case "memory":
		return memory.New()
	case "redis":
		addr, port, err := parseAddr(s.Address, 6379)
		if err != nil {
			return nil
		}
		return redis.New(s.Namespace, addr, port)
	case "zookeeper":
		addr, port, err := parseAddr(s.Address, 2181)
		if err != nil {
			return nil
		}
		return zookeeper.New(s.Namespace, addr, port)
	case "database":
		u, err := url.Parse(s.Address)
		if err != nil {
			log.Warnf("incorrect address: %s", s.Address)
			return nil
		}
		info_table := u.Query().Get("info_table")
		if info_table == "" {
			info_table = "schedule_info"
		}
		db, err := sql.Open(u.Scheme,
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				s.Username,
				s.Password,
				u.Hostname(),
				u.Port(),
				strings.Trim(u.Path, "/")),
		)
		if err != nil {
			log.Warnf("Create db instance failed: %s", err.Error())
			return nil
		}
		return database.New(
			s.Namespace,
			db,
			database.WithInfoTable(info_table),
		)
	case "etcdv2":
		addr, port, err := parseAddr(s.Address, 2379)
		if err != nil {
			return nil
		}
		return etcdv2.New(s.Namespace, []string{fmt.Sprintf("http://%s:%d", addr, port)})
	case "etcdv3":
		addr, port, err := parseAddr(s.Address, 2379)
		if err != nil {
			return nil
		}
		s, err := etcdv3.New(s.Namespace, []string{fmt.Sprintf("%s:%d", addr, port)})
		if err != nil {
			log.Warnf("Create etcdv3 instance failed: %s", err.Error())
			return nil
		}
		return s
	default:
		log.Warnf("Unknown type of storage: %s", s.Type)
		return nil
	}
}
