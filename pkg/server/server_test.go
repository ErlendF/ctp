package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This is not a good test. It shouldn't be necessary to test a function nearly devoid of actual logic.
// This test is however added as the only metric used is testcoverage.
func TestNew(t *testing.T) {
	server := New(80, &mockUserManager{}, &mockMW{})
	assert.NotNil(t, server)
}
