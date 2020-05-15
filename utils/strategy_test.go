package utils

import (
	"testing"

	"github.com/jasonjoo2010/goschedule/core/definition"
	"github.com/stretchr/testify/assert"
)

func TestToStrategyKind(t *testing.T) {
	assert.Equal(t, definition.FuncKind, ToStrategyKind("func"))
	assert.Equal(t, definition.SimpleKind, ToStrategyKind("simple"))
	assert.Equal(t, definition.UnknowKind, ToStrategyKind("simple1"))
}
