package blizzard

import (
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// ValidateBattleUser func validates a users input to *Game Overwatch
func (b *Blizzard) ValidateBattleUser(payload *models.Overwatch) error {
	logrus.Debug("ValidateBattleUser()")
	if payload == nil {
		return fmt.Errorf("no registration")
	}

	pass := false
	var region = []string{"us", "eu", "asia"}
	var platform = []string{"pc", "etc"} //("switch", "xbox", "ps4")? //TODO: validate platforms against api(?)

	// checks that region is a-ok
	for _, reg := range region {
		if payload.Region == reg {
			pass = true
		}
	}
	if !pass {
		return fmt.Errorf(models.ClientError)
	}
	pass = false

	// checks that platform is okey dokey
	for _, plat := range platform {
		if payload.Platform == plat {
			pass = true
		}
	}
	if !pass {
		return fmt.Errorf(models.ClientError)
	}

	// check that provided battle tag is correct TODO: make regex (https://us.battle.net/support/en/article/700007)?
	url := fmt.Sprintf("https://ow-api.com/v1/stats/%s/%s/%s/heroes/complete",
		payload.Platform, payload.Region, payload.BattleTag)
	resp, err := b.Get(url)
	if err != nil {
		return fmt.Errorf(models.ClientError)
	}
	defer resp.Body.Close()

	// Checks status header
	if err = models.CheckStatusCode(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

// GetBlizzardPlaytime gets playtime for PUBLIC Overwatch profiles
func (b *Blizzard) GetBlizzardPlaytime(payload *models.Overwatch) (*models.Game, error) {
	logrus.Debugf("GetBlizzardPlaytime")
	url := fmt.Sprintf("https://ow-api.com/v1/stats/%s/%s/%s/heroes/complete",
		payload.Platform, payload.Region, payload.BattleTag)

	// Tries to get a response from unreliable api
	for tries := 0; tries < 10; tries++ {
		gameStats, err := b.queryAPI(payload, url)
		if err != nil {
			if strings.Contains(err.Error(), models.NonOK) {
				return nil, err
			}
			logrus.WithError(err).Warn("Trying again")
			continue // try again....
		}

		// if it got time values from api -> return Game object
		return gameStats, nil
	}

	// returns error if no request returned valid response
	return nil, fmt.Errorf("no acceptable response from OW-api")
}

// queryAPI func returns response from the OverwatchAPI
func (b *Blizzard) queryAPI(payload *models.Overwatch, url string) (*models.Game, error) {
	var gameTime models.BlizzardResp

	// Gets statistics from the battle tag provided
	resp, err := b.Get(url)
	if err != nil {
		logrus.WithError(err).Warn("getter error")
		return nil, err
	}
	defer resp.Body.Close()

	// Checks status header
	if err = models.CheckStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	// Decodes the response to get playtime
	err = json.NewDecoder(resp.Body).Decode(&gameTime)
	if err != nil {
		logrus.WithError(err).Warn("decoder error")
		return nil, err
	}

	// Converts the returned strings to int64 (time.Duration)
	quickTime, err := nanoTime(gameTime.QuickPlayStats.CareerStats.AllHeroes.Game.TimePlayed)
	if err != nil {
		return nil, err
	}
	compTime, err := nanoTime(gameTime.CompetitiveStats.CareerStats.AllHeroes.Game.TimePlayed)
	if err != nil {
		return nil, err
	}

	// returns *Game struct
	return &models.Game{
		Name: "Overwatch",
		Time: int((quickTime + compTime).Hours())}, nil
}

// nanoTime gets a formatted time string and returns it in nanoseconds
func nanoTime(strTime string) (time.Duration, error) {
	absTime := time.Duration(0)
	parts := strings.Split(strTime, ":")

	// expected time format is "ss" or "mm:ss" or "hr:mm:ss" aka max 3 parts
	switch len(parts) {
	case 1:
		sec, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(sec) * time.Second
	case 2:
		min, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(min) * time.Minute
		sec, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(sec) * time.Second
	case 3:
		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(hour) * time.Hour
		min, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(min) * time.Minute
		sec, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}
		absTime += time.Duration(sec) * time.Second
	default:
		// returns error if parts of string exceeds 3 (hr:min:sec)
		return 0, fmt.Errorf("OW API changed the way time is encoded, got |%d| parts", len(parts))
	}

	return absTime, nil
}
