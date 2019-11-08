package valve

import (
	"ctp/pkg/models"
	"time"

	"github.com/sirupsen/logrus"
)

//Valve is a struct which contains everything necessary to handle a request related to riot
type Valve struct {
	models.Client
}

//New returns a new valve instance
func New(client models.Client) *Valve {
	return &Valve{client}
}

//GetValvePlaytime gets playtime on steam for specified game
func (v *Valve) GetValvePlaytime(game string) (*time.Duration, error) {
	logrus.Debugf("GetSteamPlaytime")
	return nil, nil
}
