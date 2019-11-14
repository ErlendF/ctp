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
	"google.golang.org/grpc/status"
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
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		return err
	}
	return nil
}

//GetUserByID gets a user from the database
func (db *Database) GetUserByID(id string) (*models.User, error) {
	doc, err := db.Collection(userCol).Doc(id).Get(db.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf(models.NotFound)
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

//GetUserByName gets a user by name
func (db *Database) GetUserByName(name string) (*models.User, error) {
	docs, err := db.Collection(userCol).Where("name", "==", name).Documents(db.ctx).GetAll()
	if err != nil {
		if status.Code(err) != codes.NotFound {
			err = fmt.Errorf(models.NotFound)
		}
		return nil, err
	}

	// checking that only one user was recieved
	switch {
	case len(docs) < 1:
		return nil, fmt.Errorf(models.NotFound)
	case len(docs) > 1:
		return nil, fmt.Errorf("Multiple users with same username")
	}

	data := docs[0].Data()
	var user models.User
	err = mapstructure.Decode(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//UpdateUser updates the relevant fields of the user
//checks for empty values
func (db *Database) UpdateUser(user *models.User) error {

	user.Name = "" // username and games are updated by dedicated functions
	user.Games = nil

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

//UpdateGames updates the games for the given user
func (db *Database) UpdateGames(user *models.User) error {
	dbUser, err := db.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// checking that each game in the database is still present in the new games array
	for _, dbGame := range dbUser.Games {
		len := len(user.Games)
		gameFound := false

		// looking for the game in user.Games
		for i := 0; i < len && !gameFound; i++ {
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

	_, err = db.Collection(userCol).Doc(user.ID).Update(db.ctx, []firestore.Update{
		{Path: "games", Value: user.Games},
	})

	return err
}

//SetUsername sets the username for the user, returns error if it is already in use
func (db *Database) SetUsername(user *models.User) error {
	dbUser, err := db.GetUserByName(user.Name)
	if err != nil && err.Error() != models.NotFound {
		return err
	}

	if dbUser != nil {
		if dbUser.Name == user.Name {
			return nil
		}

		return fmt.Errorf("Name already in use")
	}

	_, err = db.Collection(userCol).Doc(user.ID).Update(db.ctx, []firestore.Update{
		{Path: "name", Value: user.Name},
	})

	return err
}

//OverwriteUser overwrites the user specified by the user id, or creates it if it didn't exist already
func (db *Database) OverwriteUser(user *models.User) error {
	_, err := db.Collection(userCol).Doc(user.ID).Set(db.ctx, user)
	return err
}
