package blizzard

import (
	"ctp/pkg/models"
	"time"

	"github.com/sirupsen/logrus"
)

//Blizzard is a struct which contains everything necessary to handle a request related to blizzard
type Blizzard struct {
	models.Client
}

//New returns a new blizzard instance
func New(client models.Client) *Blizzard {
	return &Blizzard{client}
}

//GetBlizzardPlaytime gets playtime on steam for specified game
func (b *Blizzard) GetBlizzardPlaytime(game string) (*time.Duration, error) {
	logrus.Debugf("GetBlizzardPlaytime")
	return nil, nil
}
