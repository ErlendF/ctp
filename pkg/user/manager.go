package user

import (
	"ctp/pkg/models"
	"errors"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

// Manager is a struct which contains everything necessary
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

// GetUserByID gets the relevant info for the given user by id
func (m *Manager) GetUserByID(id string) (*models.User, error) {
	return m.db.GetUserByID(id)
}

// GetUserByName gets the relevant info for the given user by username
func (m *Manager) GetUserByName(username string) (*models.User, error) {
	user, err := m.db.GetUserByName(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SetUser updates a given user
func (m *Manager) SetUser(user *models.User) error {
	gameChanges, err := m.validateUserInfo(user)
	if err != nil {
		return err
	}

	err = m.db.UpdateUser(user)
	if err != nil {
		return err
	}

	// Updates games if there have been a change in game providers
	if gameChanges {
		return m.UpdateGames(user.ID)
	}

	return nil
}

// DeleteUser deletes the user with the given id
func (m *Manager) DeleteUser(id string, fields []string) error {
	if len(fields) == 0 {
		return m.db.DeleteUser(id)
	}

	err := m.db.DeleteFieldsFromUser(id, fields)
	if err != nil {
		return err
	}

	return m.UpdateGames(id) // Updates the games for the user, as some game providers may have been deleted
}

// UpdateGames updates all games the user has registered
func (m *Manager) UpdateGames(id string) error {
	user, err := m.db.GetUserByID(id)
	if err != nil {
		return err
	}

	var updatedGames []models.Game

	if user.Lol != nil {
		lolGame, err := m.GetLolPlaytime(user.Lol)
		if err != nil {
			return err
		}

		updatedGames = append(updatedGames, *lolGame)
	}

	if user.Overwatch != nil {
		ow, err := m.GetBlizzardPlaytime(user.Overwatch)
		if err != nil {
			return err
		}

		updatedGames = append(updatedGames, *ow)
	}

	if user.Valve != nil {
		games, err := m.GetValvePlaytime(user.Valve.ID)
		if err != nil {
			return err
		}

		updatedGames = append(games, updatedGames...)
	}

	if user.Runescape != nil {
		rs, err := m.GetRSPlaytime(user.Runescape)
		if err != nil {
			return err
		}

		updatedGames = append(updatedGames, *rs)
	}

	user.Games = updatedGames

	return m.db.UpdateGames(user)
}

// Redirect redirects the user to oauth providers
func (m *Manager) Redirect(w http.ResponseWriter, r *http.Request) {
	m.AuthRedirect(w, r)
}

// AuthCallback handles oauth callback
func (m *Manager) AuthCallback(w http.ResponseWriter, r *http.Request) (string, error) {
	id, err := m.HandleOAuth2Callback(w, r)
	if err != nil {
		return "", err
	}

	err = m.db.CreateUser(&models.User{ID: id})
	if err != nil {
		return "", err
	}

	token, err := m.GetNewToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

// UpdateRiotAPIKey updates
func (m *Manager) UpdateRiotAPIKey(key, id string) error {
	user, err := m.db.GetUserByID(id)
	if err != nil {
		return err
	}
	if !user.Admin {
		return models.ErrInvalidID // although, not entirely acurate error, it is still an invalid ID for the given action
	}

	return m.UpdateKey(key)
}

// JohanTestFunc is just a method for johan to test things :-)
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

	tmpUser.Games = append(tmpUser.Games, tmpGame, tmpGame2)
	// debug end

	tmpUser2, err := m.db.GetUserByID("117575669351657432712")
	if err != nil {
		logrus.WithError(err).Debug("Could not get user!")
		return
	}

	tmpUser.Lol = tmpUser2.Lol

	tmpUser3, err := m.db.GetUserByID("117575669351657432712")
	if err != nil {
		logrus.WithError(err).Debug("Could not get user!")
		return
	}

	logrus.Debug(tmpUser3.Lol.AccountID)
	logrus.Debug(tmpUser3.Lol.SummonerRegion)
	logrus.Debug(tmpUser3.Lol.SummonerName)

	game, err := m.GetLolPlaytime(tmpUser3.Lol)
	if err != nil {
		logrus.WithError(err).Debug("Get riot playtime oopsie!")
		return
	}

	tmpUser3.Games = append(tmpUser3.Games, *game)

	err = m.db.UpdateGames(tmpUser3)
	if err != nil {
		logrus.WithError(err).Warn("UpdateGames failed!")
		return
	}
}

// validateUserName checks if the name entered is a valid name for a user
func validateUserName(name string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9 ]{1,15}$")
	if !re.MatchString(name) {
		return errors.New("invalid username")
	}

	return nil
}

// validateUserInfo checks whether or not any information has been updated,
// and validates the updated information
func (m *Manager) validateUserInfo(user *models.User) (bool, error) {
	dbUser, err := m.db.GetUserByID(user.ID)
	if err != nil {
		return false, err
	}
	if dbUser == nil {
		dbUser = &models.User{} // makes sure that there is no invalid nilpointer dereferense
	}

	if user.Name != "" && user.Name != dbUser.Name {
		err = validateUserName(user.Name)
		if err != nil {
			return false, err
		}

		err = m.db.SetUsername(user)
		if err != nil {
			return false, err
		}
	}

	// validating each property
	// if the property is nil, or the same as stored in the database, it is considered valid
	lol, err := m.validateLol(user.Lol, dbUser.Lol)
	if err != nil {
		return false, err
	}
	ow, err := m.validateOW(user.Overwatch, dbUser.Overwatch)
	if err != nil {
		return false, err
	}
	valve, err := m.validateValve(user.Valve, dbUser.Valve)
	if err != nil {
		return false, err
	}
	rs, err := m.validateRS(user.Runescape, dbUser.Runescape)
	if err != nil {
		return false, err
	}

	changes := lol || ow || valve || rs
	return changes, nil
}

// checking that league of legends is set and that it's different from what is already stored
// if there are no changes, it doesn't need to be validated
func (m *Manager) validateLol(reg, dbReg *models.SummonerRegistration) (bool, error) {
	if reg == nil || reg == dbReg {
		return false, nil
	}

	err := m.ValidateSummoner(reg)
	if err != nil {
		return false, err
	}

	return true, nil
}

// if there are no changes, it doesn't need to be validated
func (m *Manager) validateOW(ow, dbOW *models.Overwatch) (bool, error) {
	if ow == nil || ow == dbOW {
		return false, nil
	}

	err := m.ValidateBattleUser(ow)
	if err != nil {
		return false, err
	}

	return true, nil
}

// if there are no changes, it doesn't need to be validated
func (m *Manager) validateValve(valve, dbValve *models.ValveAccount) (bool, error) {
	if valve == nil || valve == dbValve {
		return false, nil
	}

	var err error
	switch {
	case valve.ID != "":
		err = m.ValidateValveID(valve.ID)
		if err != nil {
			return false, err
		}
		valve.Username = "" // the username is not validated, nor needed. It is therefor removed
	case valve.Username != "":
		valve.ID, err = m.ValidateValveAccount(valve.Username)
		if err != nil {
			return false, err
		}
	default:
		return false, models.NewReqErrStr("invalid steam account", "invalid steam account information")
	}

	return true, nil
}

// if there are no changes, it doesn't need to be validated
func (m *Manager) validateRS(rs, dbRS *models.RunescapeAccount) (bool, error) {
	if rs == nil || rs == dbRS {
		return false, nil
	}

	err := m.ValidateRSAccount(rs)
	if err != nil {
		return false, err
	}

	return true, nil
}
