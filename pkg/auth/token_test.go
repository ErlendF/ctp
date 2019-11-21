package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// short test to check that the same id is returend by generating and validating a token
func TestTokenGenerationValidation(t *testing.T) {
	uv := &mockUserValidator{}
	auth, err := New(context.Background(), uv, 8080, "localhost", "", "", "testSecret")
	require.Nil(t, err)
	require.NotNil(t, auth)

	testID := "this is a test id"
	token, err := auth.GetNewToken(testID)
	require.Nil(t, err)
	require.NotEmpty(t, token)

	strID, err := auth.validateToken(token)
	require.Nil(t, err)
	require.Equal(t, testID, strID)
}
