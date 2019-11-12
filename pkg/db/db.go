package db

import (
	"ctp/pkg/models"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go" // Same as python's import dependency as alias.

	"google.golang.org/api/option"
)

//Database contains a firestore client and a context
type Database struct {
	*firestore.Client
	ctx context.Context
}

const userCol = "users"

//New returns a new databse
func New(key string) (*Database, error) {
	db := &Database{ctx: context.Background()}

	// We use a service account, load credentials file that you downloaded from your project's settings menu.
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

//SetUser updates a given user, or adds it if it doesn't exist already
func (db *Database) SetUser(user *models.User) error {
	_, err := db.Collection(userCol).Doc(user.ID).Set(db.ctx, user)
	return err
}

//UpdateGame updates the gametime for the given game
func (db *Database) UpdateGame(id string, tmpGame *models.Game) error {
	user, err := db.GetUser(id)
	if err != nil {
		return err
	}

	gamePresent := false

	//update game if present
	for i := range user.Games {
		if user.Games[i].Name == tmpGame.Name {
			user.Games[i].Time = tmpGame.Time
			gamePresent = true
		}
	}

	if !gamePresent {
		user.Games = append(user.Games, *tmpGame)
	}

	return db.SetUser(user)
}

//GetUser gets a user from the database
func (db *Database) GetUser(id string) (*models.User, error) {
	doc, err := db.Collection(userCol).Doc(id).Get(db.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "code = NotFound") {
			err = fmt.Errorf("NotFound")
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
