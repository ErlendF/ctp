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


//GetBlizzardPlaytime gets playtime for Public Overwatch profiles
func (b *Blizzard) GetBlizzardPlaytime(payload *models.Overwatch) (*models.Overwatch, error) {
	logrus.Debugf("GetBlizzardPlaytime")
	url := fmt.Sprintf("https://ow-api.com/v1/stats/%s/%s/%s/heroes/complete",
		payload.Platform, payload.Region, payload.ID)

	// Tries to get a response from unreliable api
	for tries := 0; tries < 10; tries++ {
		gameStats, err := b.queryAPI(payload, url)
		if err != nil {
			logrus.WithError(err).Warnf("erroror")
			if strings.Contains(err.Error(), models.NonOK) {
				return nil, err
			}
			continue	// try again....
		}
		// if it got time values from api -> continue code
		return gameStats, nil
	}

	// returns error
	return nil, fmt.Errorf("no acceptable response from OW-api")
}

// queryAPI func returns response from the OverwatchAPI
func (b *Blizzard) queryAPI(payload *models.Overwatch, url string) (*models.Overwatch ,error) {
	var gameTime models.BlizzardResp

	// Gets statistics from the BATTLE-ID provided
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

	// returns overwatch struct
	return &models.Overwatch{
		ID:       payload.ID,
		Platform: payload.Platform,
		Region:   payload.Region,
		Playtime: quickTime + compTime,
	}, nil
}

// nanoTime gets a formatted time string and returns it in nanoseconds
func nanoTime(strTime string) (time.Duration, error) {
	absTime := time.Duration(0)
	parts := strings.Split(strTime, ":")

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
		logrus.Errorf("OW API has changed the way time is encoded")
		return 0, fmt.Errorf("OW API changed the way time is encoded, got |%d| parts", len(parts))
	}

	return absTime, nil
}