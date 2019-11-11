package user

import (
	"ctp/pkg/models"

	"github.com/bxcodec/faker"
)

//Manager is a struct which contains everything necessary
type Manager struct {
	db models.Database
}

//New returns a new riot instance
func New(db models.Database) *Manager {
	return &Manager{db: db}
}

//GetUserInfo gets the relevant info for the given user
func (m *Manager) GetUserInfo(username string) (*models.UserInfo, error) {
	var info models.UserInfo

	err := faker.FakeData(&info)

	return &info, err
}

//SetUser updates a given user, or adds it if it doesn't exist already
func (m *Manager) SetUser(user *models.UserInfo) error {
	return m.db.SetUser(user)
}
