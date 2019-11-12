package riot

import (
	"ctp/pkg/models"

	"github.com/sirupsen/logrus"
)

//Riot is a struct which contains everything necessary to handle a request related to riot
type Riot struct {
	models.Client
}

//New returns a new riot instance
func New(client models.Client) *Riot {
	return &Riot{client}
}

//GetRiotPlaytime gets playtime on League of Legends
func (r *Riot) GetRiotPlaytime() (*models.Game, error) {
	logrus.Debugf("GetLolPlaytime")

	// db.UpdateGame("League", 9, "117575669351657432712")
	game := &models.Game{Name: "League", Time: 9}
	return game, nil
}
