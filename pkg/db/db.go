package db

import (
	"ctp/pkg/models"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go" // Same as python's import dependency as alias.

	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

//CreateUser creates a user
func (db *Database) CreateUser(user *models.User) error {
	_, err := db.Collection(userCol).Doc(user.ID).Create(db.ctx, user)
	return err
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

//UpdateUser updates the relevant fields of the user
//checks for empty values
func (db *Database) UpdateUser(user *models.User) error {
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

//GetUser gets a user from the database
func (db *Database) GetUser(id string) (*models.User, error) {
	doc, err := db.Collection(userCol).Doc(id).Get(db.ctx)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
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
