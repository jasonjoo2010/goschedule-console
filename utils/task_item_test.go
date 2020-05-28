package utils

import (
	"testing"

	"github.com/jasonjoo2010/goschedule/core/definition"
	"github.com/stretchr/testify/assert"
)

func TestParseTaskItems(t *testing.T) {
	var list []*definition.TaskItem

	list = ParseTaskItems("")
	assert.Equal(t, 0, len(list))

	list = ParseTaskItems("a,b,c,d")
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "a", list[0].Id)
	assert.Empty(t, list[0].Parameter)
	assert.Equal(t, "d", list[3].Id)
	assert.Empty(t, list[3].Parameter)

	list = ParseTaskItems("a{asdf},b,c:,d")
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "a{asdf}", list[0].Id)
	assert.Empty(t, list[0].Parameter)
	assert.Equal(t, "d", list[3].Id)
	assert.Empty(t, list[3].Parameter)

	list = ParseTaskItems("a:{asdf},b,c:{a=1,b=2},d")
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "a", list[0].Id)
	assert.Equal(t, "asdf", list[0].Parameter)
	assert.Equal(t, "c", list[2].Id)
	assert.Equal(t, "a=1,b=2", list[2].Parameter)
}
