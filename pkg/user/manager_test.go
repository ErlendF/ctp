package user

import (
	"ctp/pkg/models"
	"errors"
	"net/http"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	err  error
	user *models.User
}

func (m *mockDB) CreateUser(user *models.User) error              { return m.err }
func (m *mockDB) GetUserByID(id string) (*models.User, error)     { return m.user, m.err }
func (m *mockDB) GetUserByName(name string) (*models.User, error) { return m.user, m.err }
func (m *mockDB) UpdateUser(user *models.User) error              { return m.err }
func (m *mockDB) UpdateGames(user *models.User) error             { return m.err }
func (m *mockDB) SetUsername(user *models.User) error             { return m.err }
func (m *mockDB) OverwriteUser(user *models.User) error           { return m.err }
func (m *mockDB) DeleteUser(id string) error                      { return m.err }

type mockOrganizer struct {
	valve   []models.Game
	valveID string
	lol     *models.Game
	rs      *models.Game
	ow      *models.Game
	id      string
	token   string
	err     error
}

func (m *mockOrganizer) ValidateValveAccount(username string) (string, error) { return m.valveID, m.err }
func (m *mockOrganizer) GetValvePlaytime(id string) ([]models.Game, error)    { return m.valve, m.err }
func (m *mockOrganizer) GetRiotPlaytime(reg *models.SummonerRegistration) (*models.Game, error) {
	return m.lol, m.err
}
func (m *mockOrganizer) ValidateSummoner(reg *models.SummonerRegistration) error { return m.err }
func (m *mockOrganizer) GetRSPlaytime(username string) (*models.Game, error)     { return m.rs, m.err }
func (m *mockOrganizer) ValidateRSAccount(username string) error                 { return m.err }
func (m *mockOrganizer) GetBlizzardPlaytime(*models.Overwatch) (*models.Game, error) {
	return m.ow, m.err
}
func (m *mockOrganizer) ValidateBattleUser(overwatch *models.Overwatch) error { return m.err }
func (m *mockOrganizer) GetNewToken(id string) (string, error)                { return m.token, m.err }
func (m *mockOrganizer) AuthRedirect(w http.ResponseWriter, r *http.Request)  {}
func (m *mockOrganizer) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error) {
	return m.id, m.err
}

func TestSetUser(t *testing.T) {
	var cases = []struct {
		name        string
		orgErr      error
		dbErr       error
		expectedErr error
	}{
		{"Test ok", nil, nil, nil},
		{"Test orgErr", errors.New("test"), nil, errors.New("test")},
		{"Test orgErr", nil, errors.New("test"), errors.New("test")},
	}

	db := &mockDB{}
	org := &mockOrganizer{}
	um := New(db, org)

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var user *models.User
			err := faker.FakeData(&user)
			assert.NoError(t, err)
			user.Name = "testuser123"

			err = faker.FakeData(&db.user)
			assert.NoError(t, err)
			db.err = tc.dbErr
			fakeOrg(t, org, tc.orgErr)

			err = um.SetUser(user)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func fakeOrg(t *testing.T, org *mockOrganizer, orgErr error) {
	err := faker.FakeData(&org.lol)
	assert.NoError(t, err)
	err = faker.FakeData(&org.ow)
	assert.NoError(t, err)
	err = faker.FakeData(&org.rs)
	assert.NoError(t, err)
	err = faker.FakeData(&org.valve)
	assert.NoError(t, err)
	err = faker.FakeData(&org.valveID)
	assert.NoError(t, err)
	err = faker.FakeData(&org.id)
	assert.NoError(t, err)
	err = faker.FakeData(&org.token)
	assert.NoError(t, err)
	org.err = orgErr
}
