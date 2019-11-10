package db

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go" // Same as python's import dependency as alias.

	"google.golang.org/api/option"
)

//Database contains a firestore client and a context
type Database struct {
	*firestore.Client
	ctx context.Context
}

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
