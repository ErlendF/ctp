package riot

import (
	"ctp/pkg/models"
	"time"

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
func (r *Riot) GetRiotPlaytime() (*time.Duration, error) {
	logrus.Debugf("GetLolPlaytime")

	


	return nil, nil
}
