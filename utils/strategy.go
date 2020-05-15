package utils

import (
	"github.com/jasonjoo2010/goschedule/core/definition"
)

// ToStrategyKind returns kind of strategy and 0 for error
func ToStrategyKind(str string) definition.StrategyKind {
	switch str {
	case "simple":
		return definition.SimpleKind
	case "func":
		return definition.FuncKind
	case "task":
		return definition.TaskKind
	}
	return 0
}
