package auth

import (
	"context"
	"ctp/pkg/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	t          *testing.T
	expectedID string
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(models.CtxKey("id"))
	idStr, ok := id.(string)

	if !ok {
		assert.Fail(m.t, "invalid id")
	}

	assert.Equal(m.t, m.expectedID, idStr)
}

type mockUserValidator struct {
	resp bool
	err  error
}

func (m *mockUserValidator) IsUser(id string) (bool, error) { return m.resp, m.err }

func TestAuthMiddleware(t *testing.T) {
	var cases = []struct {
		name           string
		err            error
		uvResp         bool
		provideToken   bool
		expectedStatus int
	}{
		{"Test ok", nil, true, true, http.StatusOK},
		{"Test no token provided", nil, true, false, http.StatusForbidden},
		{"Test unexpected error", errors.New("test"), true, false, http.StatusForbidden},
		{"Test invalid user", nil, true, false, http.StatusForbidden},
	}

	uv := &mockUserValidator{}
	auth, err := New(context.Background(), uv, 8080, "localhost", "", "", "testSecret")
	require.NoError(t, err)

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Initializing mock structs with random data
			var id string
			err = faker.FakeData(&id)
			require.Nil(t, err)
			h := &mockHandler{t: t, expectedID: id}
			mw := auth.Auth(h)

			uv.err = tc.err
			uv.resp = tc.uvResp

			// Making and serving request
			req, err := http.NewRequest("GET", "test", nil) // both method and url is handled by the router
			require.Nil(t, err)

			token, err := auth.GetNewToken(id)
			require.Nil(t, err)

			if tc.provideToken {
				req.Header.Set("Authorization", token)
			}

			w := httptest.NewRecorder()
			mw.ServeHTTP(w, req)
			resp := w.Result()
			if resp.Body != nil {
				defer resp.Body.Close()
			}

			// Body should only be parsed if expected to succeed, and it actually succeeded
			//  should assert.Equal regardless of expected status
			if !assert.Equal(t, tc.expectedStatus, resp.StatusCode) || tc.expectedStatus != http.StatusOK {
				return
			}
		})
	}
}
