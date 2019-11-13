package user

import (
	"ctp/pkg/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//Manager is a struct which contains everything necessary
type Manager struct {
	models.Organizer
	db models.Database
}

// New returns a new user manager instance.
// The manager takes a db and organizer. It embedds the organizer to simplify calls
// Organizer is used to simplify the passing of all interfaces to the handler.
func New(db models.Database, organizer models.Organizer) *Manager {
	m := &Manager{db: db}
	m.Organizer = organizer
	return m
}

//GetUser gets the relevant info for the given user
func (m *Manager) GetUser(id string) (*models.User, error) {
	return m.db.GetUser(id)
}

//SetUser updates a given user
func (m *Manager) SetUser(user *models.User) error {
	var err error
	if user.Name != "" {
		err = validateUserName(user.Name)
		if err != nil {
			return err
		}
	}

	if user.Lol != nil {
		user.Lol, err = m.ValidateSummoner(user.Lol)
		if err != nil {
			return err
		}
	}

	//TODO: validate steam and other ids or registrations

	return m.db.UpdateUser(user)
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
	m.AuthRedirect(w, r)
}

//AuthCallback handles oauth callback
func (m *Manager) AuthCallback(w http.ResponseWriter, r *http.Request) (string, error) {
	id, err := m.HandleOAuth2Callback(w, r)

	err = m.db.CreateUser(&models.User{ID: id})
	if err != nil {
		if grpc.Code(err) != codes.AlreadyExists {
			return "", err
		}
		err = nil
	}

	token, err := m.GetNewToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

//RegisterLeague registeres League of Legends for a given user
func (m *Manager) RegisterLeague(id string, reg *models.SummonerRegistration) error {
	reg, err := m.ValidateSummoner(reg)
	if err != nil {
		return err
	}

	user := &models.User{ID: id, Lol: reg}

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
		Valve:         "76561198075109466",
	}

	tmpUser.Games = append(tmpUser.Games, tmpGame)
	tmpUser.Games = append(tmpUser.Games, tmpGame2)
	//debug end

	tmpUser2, err := m.db.GetUser("117575669351657432712")
	if err != nil {
		logrus.WithError(err).Debugf("Could not get user!")
		return
	}
	tmpUser.Lol = tmpUser2.Lol

	err = m.db.SetUser(&tmpUser)
	if err != nil {
		logrus.WithError(err).Debugf("Test failed!")
	}

	tmpUser3, err := m.db.GetUser("117575669351657432712")
	if err != nil {
		logrus.WithError(err).Debugf("Could not get user!")
		return
	}

	logrus.Debugf(tmpUser3.Lol.AccountID)
	logrus.Debugf(tmpUser3.Lol.SummonerRegion)
	logrus.Debugf(tmpUser3.Lol.SummonerName)

	game, err := m.GetRiotPlaytime(tmpUser3.Lol)
	if err != nil {
		logrus.WithError(err).Debugf("Get riot playtime oopsie!")
		return
	}

	err = m.db.UpdateGame("117575669351657432712", game)
	if err != nil {
		logrus.WithError(err).Warnf("Update game failed!")
		return
	}

	games, err := m.GetValvePlaytime(tmpUser3.Valve)
	if err != nil {
		logrus.WithError(err).Warnf("Valve playtime failed!")
		return
	}

	for _, game := range games {
		logrus.Debugf(game.Name)
		err = m.db.UpdateGame("117575669351657432712", &game)
		if err != nil {
			logrus.WithError(err).Warnf("Update game failed!")
			return
		}
	}
}

//validateUserName checks if the name entered is a valid name for a user
func validateUserName(name string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9 ]{1,15}$")
	if !re.MatchString(name) {
		return fmt.Errorf("Invalid username")
	}
	return nil
}
