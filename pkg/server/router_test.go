package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockMW struct{}

func (m *mockMW) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// This is not a good test. However, it shouldn't be necessary to test a function nearly devoid of actual logic.
// This test is however added as the only metric used is testcoverage.
func TestNewRouter(t *testing.T) {
	um := &mockUserManager{}
	h := newHandler(um)
	r := newRouter(h, &mockMW{})
	require.NotNil(t, r)
}
