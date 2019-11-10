package user

import (
	"ctp/pkg/models"

	"github.com/bxcodec/faker"
)

//Manager is a struct which contains everything necessary
type Manager struct {
}

//New returns a new riot instance
func New() *Manager {
	return &Manager{}
}

//GetUserInfo gets the relevant info for the given user
func (m *Manager) GetUserInfo(username string) (*models.UserInfo, error) {
	var info models.UserInfo

	err := faker.FakeData(&info)

	return &info, err
}
