package server

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ctp/pkg/models"

	"github.com/bxcodec/faker"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserManager struct {
	user     *models.User
	response string
	err      error
}

func (m *mockUserManager) GetUserByID(id string) (*models.User, error)         { return m.user, m.err }
func (m *mockUserManager) GetUserByName(username string) (*models.User, error) { return m.user, m.err }
func (m *mockUserManager) SetUser(user *models.User) error                     { return m.err }
func (m *mockUserManager) DeleteUser(id string) error                          { return m.err }
func (m *mockUserManager) UpdateGames(id string) error                         { return m.err }
func (m *mockUserManager) Redirect(w http.ResponseWriter, r *http.Request)     {}
func (m *mockUserManager) JohanTestFunc()                                      {}
func (m *mockUserManager) AuthCallback(w http.ResponseWriter, r *http.Request) (string, error) {
	return m.response, m.err
}

func TestHandler(t *testing.T) {
	var cases = []struct {
		name           string
		err            error
		url            string
		reqBody        string
		method         string
		expectedStatus int
	}{
		{"Test ok return for GET /user", nil, "/api/v1/user", "", http.MethodGet, http.StatusOK},
		{"Test err not found for GET /user", models.ErrNotFound, "/api/v1/user", "", http.MethodGet, http.StatusNotFound},
		{"Test invalid id for GET /user", models.ErrInvalidID, "/api/v1/user", "", http.MethodGet, http.StatusForbidden},
		{"Test not found for GET /user", models.ErrNotFound, "/api/v1/user", "", http.MethodGet, http.StatusNotFound},
		{"Test request error for GET /user", models.NewReqErrStr("test", "resp"), "/api/v1/user", "", http.MethodGet, http.StatusBadRequest},
		{"Test API error for GET /user", models.NewAPIErr(errors.New("test"), "Test"), "/api/v1/user", "", http.MethodGet, http.StatusBadGateway},
		{"Test unexpected error for GET /user", errors.New("test"), "/api/v1/user", "", http.MethodGet, http.StatusInternalServerError},
		{"Test ok return for POST /user", nil, "/api/v1/user",
			`{
			"username": "newUsername",
			"lol": {
				"summonerName": "LOPER",
				"summonerRegion": "EUW1"
			},
			"valve": {
				"username": "test"
			},
			"overwatch": {
				"battleTag": "Onijuan-2670",
				"platform": "pc",
				"region": "eu"
			}
		}`, http.MethodPost, http.StatusOK},
		{"Test invalid json request body for POST /user", nil, "/api/v1/user", `{ this is an invalid request body }`, http.MethodPost, http.StatusBadRequest},
		//{"Test invalid request body for POST /user", nil, "/api/v1/user", `{ "test": "this is valid json, but invalid request body" }`, http.MethodPost, http.StatusBadRequest}, - Requires changes to handler, fields are just ignored
		{"Test ok return for DELETE /user", nil, "/api/v1/user", "", http.MethodDelete, http.StatusOK},
		{"Test ok return for POST /updategames", nil, "/api/v1/updategames", "", http.MethodPost, http.StatusOK},
		{"Test ok return for GET /authcallback", nil, "/api/v1/authcallback", "", http.MethodGet, http.StatusOK},
		{"Test ok return for GET /login", nil, "/api/v1/login", "", http.MethodGet, http.StatusOK}, // as the redirection is never called (due to mocking), the function returns statusOK
		{"Test ok return for GET /user/{username}", nil, "/api/v1/user/test", "", http.MethodGet, http.StatusOK},
		{"Test invalid username GET /user/{username}", nil, "/api/v1/user/012345678901234567890", "", http.MethodGet, http.StatusNotFound},
	}

	um := &mockUserManager{}
	h := newHandler(um)
	assert.NotNil(t, h)

	r := mockRouter(h)

	k := ctxKey("id")

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userResp := &models.User{}
			// Initializing mock structs with random data
			err := faker.FakeData(&um.user)
			require.Nil(t, err)
			err = faker.FakeData(&um.response)
			require.Nil(t, err)

			um.err = tc.err

			// Making and serving request
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.reqBody))
			require.Nil(t, err)

			ctx := context.WithValue(req.Context(), k, "12345")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			resp := w.Result()
			if resp.Body != nil {
				defer resp.Body.Close()
			}

			// Body should only be parsed if expected to succeed, and it actually succeeded
			//  should assert.Equal regardless of expected status
			if !assert.Equal(t, tc.expectedStatus, resp.StatusCode) || tc.expectedStatus != http.StatusOK {
				return
			}

			// Decoding returned data and comparing with data from mock structs
			if strings.Contains(tc.url, "/api/v1/user") && tc.method == http.MethodGet {
				err = json.NewDecoder(resp.Body).Decode(&userResp)
				assert.Nil(t, err)
				removeIgnoredOutput(um.user, tc.url)
				assert.Equal(t, um.user, userResp)
			} else if tc.url == "/api/v1/authcallback" {
				body, err := ioutil.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.Equal(t, um.response, string(body))
			} else if tc.url != "/api/v1/login" {
				body, err := ioutil.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.Equal(t, "Success", string(body))
			}
		})
	}
}

// need to make a new router to avoid testing the middleware aswell
func mockRouter(h *handler) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.notFound)

	get := r.PathPrefix("/api/v1").Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/", h.testHandler).Name("root")
	get.HandleFunc("/login", h.login).Name("login")
	get.HandleFunc("/authcallback", h.authCallbackHandler).Name("authCallback")
	get.HandleFunc("/user/{username:[a-zA-Z0-9 ]{1,15}}", h.getPublicUser).Name("getPublicUser")

	auth := r.PathPrefix("/api/v1/").Subrouter()
	auth.HandleFunc("/user", h.getUser).Methods(http.MethodGet).Name("getUser")
	auth.HandleFunc("/user", h.updateUser).Methods(http.MethodPost).Name("updateUser")
	auth.HandleFunc("/user", h.deleteUser).Methods(http.MethodDelete).Name("deleteUser")
	auth.HandleFunc("/updategames", h.updateGames).Methods(http.MethodPost).Name("updateGames")

	return r
}

// userID and valveID is not returned to the user. valveID is only used to differentiate the games internaly
func removeIgnoredOutput(user *models.User, url string) {
	user.ID = ""
	if strings.Contains(url, "/api/v1/user/") { //it's a test for /api/v1/user/{username}
		user.Public = false
	}
	for i := range user.Games {
		user.Games[i].ValveID = 0
	}
}
