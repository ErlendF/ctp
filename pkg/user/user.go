package user

import (
	"ctp/pkg/models"
	"net/http"

	"github.com/sirupsen/logrus"
)

//Manager is a struct which contains everything necessary
type Manager struct {
	o  models.Organizer
	db models.Database
}

//New returns a new user manager instance
// Organizer is used to simplify the passing of all interfaces to the handler.
func New(db models.Database, organizer models.Organizer) *Manager {
	m := &Manager{db: db, o: organizer}
	return m
}

//GetUser gets the relevant info for the given user
func (m *Manager) GetUser(id string) (*models.User, error) {
	return m.db.GetUser(id)
}

//SetUser updates a given user
func (m *Manager) SetUser(user *models.User) error {
	return m.db.SetUser(user)
}

//UpdateGame updates the gametime for the given game
func (m *Manager) UpdateGame(id string, game *models.Game) error {
	return m.db.UpdateGame(id, game)
}

//UpdateAllGames updates all games the user has registered
func (m *Manager) UpdateAllGames(id string) error {
	return nil
}

//Redirect redirects the user to oauth providers
func (m *Manager) Redirect(w http.ResponseWriter, r *http.Request) {
	m.o.Redirect(w, r)
}

//AuthCallback handles oauth callback
func (m *Manager) AuthCallback(w http.ResponseWriter, r *http.Request) (string, error) {
	id, err := m.o.HandleOAuth2Callback(w, r)

	err = m.SetUser(&models.User{ID: id})
	if err != nil {
		return "", err
	}

	token, err := m.o.GetNewToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

//RegisterLeague registeres League of Legends for a given user
func (m *Manager) RegisterLeague(id string, reg *models.SummonerRegistration) error {
	reg, err := m.o.ValidateSummoner(reg)
	if err != nil {
		return err
	}

	user := &models.User{ID: id}

	user.Lol = *reg

	return m.db.UpdateUser(user)
}

//JohanTestFunc is just a method for johan to test things :-)
func (m *Manager) JohanTestFunc() {
	tmpGame := models.Game{
		Name: "League",
		Time: 12,
	}

	tmpGame2 := models.Game{
		Name: "RocketLeage",
		Time: 112,
	}

	tmpUser := models.User{
		ID:            "117575669351657432712",
		Name:          "Johan",
		TotalGameTime: 12,
		Games:         nil,
	}

	tmpUser.Games = append(tmpUser.Games, tmpGame)
	tmpUser.Games = append(tmpUser.Games, tmpGame2)
	//debug end

	err := m.db.SetUser(&tmpUser)
	if err != nil {
		logrus.WithError(err).Debugf("Test failed!")
	}


	tmpUser2, _ := m.db.GetUser("117575669351657432712")
	game, _ := m.o.GetRiotPlaytime(tmpUser2.Lol)

	err = m.UpdateGame("117575669351657432712", game)
	if err != nil {
		logrus.WithError(err).Warnf("Update game failed!")
	}
}
