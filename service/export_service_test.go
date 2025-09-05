package service

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestNumberToColumn(t *testing.T) {
	assert.Equal(t, NumberToColumn(0), "A")
	assert.Equal(t, NumberToColumn(1), "B")
	assert.Equal(t, NumberToColumn(2), "C")
	assert.Equal(t, NumberToColumn(25), "Z")
	assert.Equal(t, NumberToColumn(26), "AA")
	assert.Equal(t, NumberToColumn(27), "AB")
	assert.Equal(t, NumberToColumn(28), "AC")
}
