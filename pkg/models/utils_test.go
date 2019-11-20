package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	var cases = []struct {
		name     string
		expected bool
		item     string
		arr      []string
	}{
		{"Test item is contained", true, "test", []string{"test", "test1", "test2"}},
		{"Test item is not contained", false, "test", []string{"test1", "test2"}},
	}

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Contains(tc.arr, tc.item)
			assert.Equal(t, tc.expected, result)
		})
	}
}
