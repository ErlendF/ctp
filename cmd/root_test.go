package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestFlags(t *testing.T) {
	assert.Equal(t, config.port, 80)
	assert.Equal(t, config.verbose, false)
}

