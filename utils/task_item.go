// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package utils

import (
	"strings"

	"github.com/jasonjoo2010/goschedule/core/definition"
)

// a:{param1},bbb:{param2}
// 01222222230000122222223  (state transmition)
// 0 -> id segment, wait for ':' or ','(reset)
// 1 -> wait for '{' or ','(reset)
// 2 -> id segment, wait for '}', no reset
// 3 -> wait for ',' to reset
type state struct {
	result []definition.TaskItem
	step   int
	id     strings.Builder
	param  strings.Builder
}

func (s *state) Reset() {
	s.step = 0
	if s.id.Len() == 0 {
		if s.param.Len() > 0 {
			s.param.Reset()
		}
		return
	}
	item := definition.TaskItem{Id: s.id.String()}
	if s.param.Len() > 0 {
		item.Parameter = s.param.String()
		s.param.Reset()
	}
	s.id.Reset()
	s.result = append(s.result, item)
}

func (s *state) Process(b byte) {
	switch {
	case s.step == 0 && b == ':':
		s.step++
	case s.step == 1 && b == '{':
		s.step++
	case s.step == 2 && b == '}':
		s.step++
	case s.step == 3 && b == ',':
		// normal reset
		s.Reset()
	case (s.step == 0 || s.step == 1) && b == ',':
		// abnormal reset
		s.Reset()
	case s.step == 0:
		s.id.WriteByte(b)
	case s.step == 2:
		s.param.WriteByte(b)
	}
}

func ParseTaskItems(str string) []definition.TaskItem {
	if str == "" {
		return make([]definition.TaskItem, 0)
	}
	// use state machine
	state := &state{
		result: make([]definition.TaskItem, 0, 2),
	}
	for i := 0; i < len(str); i++ {
		state.Process(str[i])
	}
	state.Reset()
	return state.result
}
