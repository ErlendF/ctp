package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	assert.Equal(t, config.port, 80)
	assert.Equal(t, config.verbose, false)
}
