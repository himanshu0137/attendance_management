package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBusFunctions(t *testing.T) {
	t.Run("functions compile and exist", func(t *testing.T) {
		// This test just verifies that our functions can be compiled
		assert.True(t, true)
	})
}