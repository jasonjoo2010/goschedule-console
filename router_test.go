package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func judge(t *testing.T, result interface{}, expected int64) {
	assert.IsType(t, expected, result)
	assert.Equal(t, expected, result)
}

func judgef(t *testing.T, result interface{}, expected float64) {
	assert.IsType(t, expected, result)
	assert.Equal(t, expected, result)
}

func TestAdd(t *testing.T) {
	judge(t, add(int(1), int(2)), 3)
	judge(t, add(int8(1), int8(2)), 3)
	judge(t, add(int8(1), int64(2)), 3)
	judgef(t, add(int8(1), 2.0), 3)
}
