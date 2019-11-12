package blizzard

import (
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

//Blizzard is a struct which contains everything necessary to handle a request related to blizzard
type Blizzard struct {
	models.Getter
}

//New returns a new blizzard instance
func New(getter models.Getter) *Blizzard {
	return &Blizzard{getter}
}

//GetBlizzardPlaytime gets playtime for Public Overwatch profiles
func (b *Blizzard) GetBlizzardPlaytime(platform, region, ID string) (*time.Duration, error) {
	logrus.Debugf("GetBlizzardPlaytime")

	// Gets statistics from the BATTLE-ID provided
	resp, err := b.Get(fmt.Sprintf("https://ow-api.com/v1/stats/%s/%s/%s/heroes/complete", platform, region, ID))
	if err != nil {
		logrus.Errorf("Getter error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Decodes the response to get playtime.
	var gameTime models.BlizzardResp
	err = json.NewDecoder(resp.Body).Decode(&gameTime)
	if err != nil {
		logrus.Errorf("Decoder error: %v", err)
		return nil, err
	}

	// Converts the returned strings to int64 (time.Duration)
	quickTime, err := nanoTime(gameTime.QuickPlayStats.AllHeroes.Game.TimePlayed)
	if err != nil {
		logrus.Errorf("Error splitting time: %v", err)
		return nil, err
	}
	compTime, err := nanoTime(gameTime.QuickPlayStats.AllHeroes.Game.TimePlayed)
	if err != nil {
		logrus.Errorf("Error splitting time: %v", err)
		return nil, err
	}

	xd := quickTime + compTime
	return &xd, nil
}

// nanoTime gets a formatted time string and returns it in nanoseconds
func nanoTime(time string) (time.Duration, error) {



	return 0, nil
}