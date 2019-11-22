package blizzard

import (
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Blizzard is a struct which contains everything necessary to handle a request related to blizzard
type Blizzard struct {
	models.Getter
}

// blizzardResp struct retrieves only time played in Overwatch
type blizzardResp struct {
	CompetitiveStats struct {
		CareerStats careerStats `json:"careerStats"`
	} `json:"competitiveStats"`
	QuickPlayStats struct {
		CareerStats careerStats `json:"careerStats"`
	} `json:"quickPlayStats"`
}

// careerStats includes time played for all heroes by game
type careerStats struct {
	AllHeroes struct {
		Game struct {
			TimePlayed string `json:"timePlayed"`
		} `json:"game"`
	} `json:"allHeroes"`
}

func New(getter models.Getter) *Blizzard {
	return &Blizzard{getter}
}

// errInvalidTimePlayed is used to indicate to try the request again
var errInvalidTimePlayed = errors.New("invalid time played in response")

// ValidateBattleUser func validates a users input to *Game Overwatch
func (b *Blizzard) ValidateBattleUser(payload *models.Overwatch) error {
	logrus.Debug("ValidateBattleUser()")

	if payload == nil {
		return errors.New("no payload to ValidateBattleUser")
	}

	var regions = []string{"us", "eu", "asia"}

	var platforms = []string{"pc", "switch", "xbox", "ps4"}

	if !models.Contains(regions, payload.Region) {
		return models.NewReqErrStr("invalid Overwatch region", "invalid region for Overwatch account")
	}

	if !models.Contains(platforms, payload.Platform) {
		return models.NewReqErrStr("invalid Overwatch platform", "invalid platform for Overwatch account")
	}

	// check that provided battle tag is correct
	url := fmt.Sprintf("https://ow-api.com/v1/stats/%s/%s/%s/heroes/complete",
		payload.Platform, payload.Region, payload.BattleTag)
	resp, err := b.Get(url)

	if err != nil {
		return models.NewAPIErr(err, "Blizzard")
	}
	defer resp.Body.Close()

	// Checks status header
	if err := models.AccValStatusCode(resp.StatusCode, "Blizzard", "invalid Blizzard battle tag, platform or region"); err != nil {
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
		gameStats, err := b.queryAPI(url)
		if err != nil {
			if !errors.Is(err, errInvalidTimePlayed) {
				return nil, models.NewAPIErr(err, "Blizzard")
			}

			continue // try again....
		}

		// got valid time values, thus returning Game object
		return gameStats, nil
	}

	// returns error if no request returned valid response
	return nil, models.NewAPIErr(errors.New("no acceptable response from OW-api"), "Blizzard")
}

// queryAPI func returns response from the OverwatchAPI
func (b *Blizzard) queryAPI(url string) (*models.Game, error) {
	var gameTime blizzardResp

	// Gets statistics from the battle tag provided
	resp, err := b.Get(url)
	if err != nil {
		return nil, models.NewAPIErr(err, "Blizzard")
	}
	defer resp.Body.Close()

	// Checks status code
	err = models.CheckStatusCode(resp.StatusCode, "Blizzard", "invalid Blizzard battle tag or region")
	if err != nil {
		return nil, err
	}

	// Decodes the response to get playtime
	err = json.NewDecoder(resp.Body).Decode(&gameTime)
	if err != nil {
		return nil, models.NewAPIErr(err, "Blizzard")
	}

	if gameTime.QuickPlayStats.CareerStats.AllHeroes.Game.TimePlayed == "" ||
		gameTime.CompetitiveStats.CareerStats.AllHeroes.Game.TimePlayed == "" {
		return nil, errInvalidTimePlayed
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
