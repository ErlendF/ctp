package user

import (
	"ctp/pkg/models"
)

//Manager is a struct which contains everything necessary
type Manager struct {
	db models.Database
}

//New returns a new riot instance
func New(db models.Database) *Manager {
	return &Manager{db: db}
}

//GetUser gets the relevant info for the given user
func (m *Manager) GetUser(id string) (*models.User, error) {
	return m.db.GetUser(id)
}

//SetUser updates a given user
func (m *Manager) SetUser(user *models.User) error {
	return m.db.SetUser(user)
}

// //AddUser adds a new user
// func (m *Manager) SetUser(user *models.User) (string, error) {
// 	err := m.db.SetUser(user)
// 	if err != nil {
// 		return "", err
// 	}
// 	token, err := m.
// }
