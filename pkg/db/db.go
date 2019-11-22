package db

import (
	"ctp/pkg/models"
	"errors"
	"strings"

	"sort"

	"github.com/sirupsen/logrus"

	"cloud.google.com/go/firestore"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go" // Same as python's import dependency as alias.

	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Database contains a firestore client and a context
type Database struct {
	*firestore.Client
	ctx context.Context
}

const userCol = "users"

var deletableFields = [...]string{"name", "games", "lol", "valve", "overwatch", "runescape", "games"}

// New returns a new databse containing context and a firestore client
func New(key string) (*Database, error) {
	db := &Database{ctx: context.Background()}

	// We use a service account. The key location defaults to "./fbkey.json", but can be configured by the "-f" flag
	opt := option.WithCredentialsFile(key)
	app, err := firebase.NewApp(db.ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	db.Client, err = app.Firestore(db.ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUser creates a user
func (db *Database) CreateUser(user *models.User) error {
	_, err := db.Collection(userCol).Doc(user.ID).Create(db.ctx, user)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}

	return nil
}

// GetUserByID gets a user from the database
func (db *Database) GetUserByID(id string) (*models.User, error) {
	doc, err := db.Collection(userCol).Doc(id).Get(db.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	data := doc.Data()

	var user models.User

	err = mapstructure.Decode(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByName gets a user by name
func (db *Database) GetUserByName(name string) (*models.User, error) {
	docs, err := db.Collection(userCol).Where("name", "==", name).Where("public", "==", true).Documents(db.ctx).GetAll()
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	// checking that only one user was received
	switch {
	case len(docs) < 1:
		return nil, models.ErrNotFound
	case len(docs) > 1:
		return nil, errors.New("multiple users with same username")
	}

	data := docs[0].Data()

	var user models.User

	err = mapstructure.Decode(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates the relevant fields of the user
// checks for empty values
func (db *Database) UpdateUser(user *models.User) error {
	user.Name = "" // username and games are updated by dedicated functions
	user.Games = nil
	user.Admin = false // admin has to be set manually

	s := structs.New(user)
	m := make(map[string]interface{})

	for _, f := range s.Fields() {
		if !f.IsZero() {
			m[f.Tag("firestore")] = f.Value()
		}
	}

	_, err := db.Collection(userCol).Doc(user.ID).Set(db.ctx, m, firestore.MergeAll)

	return err
}

// UpdateGames updates the games for the given user
func (db *Database) UpdateGames(user *models.User) error {
	dbUser, err := db.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// checking that each game in the database is still present in the new games array
	for _, dbGame := range dbUser.Games {
		gamesLen := len(user.Games)
		gameFound := false

		// looking for the game in user.Games
		for i := 0; i < gamesLen && !gameFound; i++ {
			if dbGame.Name == user.Games[i].Name && dbGame.ValveID == user.Games[i].ValveID {
				gameFound = true
				break
			}
		}

		// if the game was not present, adding it
		if !gameFound {
			user.Games = append(user.Games, dbGame)
		}
	}

	sort.Slice(user.Games, func(i, j int) bool {
		return user.Games[i].Time > user.Games[j].Time
	})

	_, err = db.Collection(userCol).Doc(user.ID).Update(db.ctx, []firestore.Update{
		{Path: "games", Value: user.Games},
	})

	if err != nil {
		return err
	}

	err = db.UpdateTotalGameTime(user.ID)

	return err
}

// UpdateTotalGameTime updates the totalgametime for the given user
func (db *Database) UpdateTotalGameTime(id string) error {
	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}

	logrus.Debugf("UpdateTotalGameTime")

	totalGameTime := 0

	for _, game := range user.Games {
		totalGameTime += game.Time
	}

	_, err = db.Collection(userCol).Doc(id).Update(db.ctx, []firestore.Update{
		{Path: "totalGameTime", Value: totalGameTime},
	})

	return err
}

// SetUsername sets the username for the user, returns error if it is already in use
func (db *Database) SetUsername(user *models.User) error {
	user.Name = strings.ToLower(user.Name)
	dbUser, err := db.GetUserByName(user.Name)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	if dbUser != nil {
		if dbUser.Name == user.Name {
			return nil
		}

		return errors.New("name already in use")
	}

	_, err = db.Collection(userCol).Doc(user.ID).Update(db.ctx, []firestore.Update{
		{Path: "name", Value: user.Name},
	})

	return err
}

// DeleteUser deletes a user from the database
func (db *Database) DeleteUser(id string) error {
	_, err := db.Collection(userCol).Doc(id).Delete(db.ctx)
	return err
}

func (db *Database) DeleteFieldsFromUser(id string, fields []string) error {
	if len(fields) > len(deletableFields) {
		return models.NewReqErrStr("too many fields to delete", "invalid request body: too many specified fields to delete")
	}
	m := make(map[string]interface{})

	for _, f := range deletableFields {
		if models.Contains(fields, f) {
			m[f] = firestore.Delete
		}
	}

	_, err := db.Collection(userCol).Doc(id).Set(db.ctx, m, firestore.MergeAll)
	return err
}

// IsUser checks wether or not the provided user exisits in the database
func (db *Database) IsUser(id string) (bool, error) {
	_, err := db.Collection(userCol).Doc(id).Get(db.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
