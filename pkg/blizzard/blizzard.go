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
	compTime, err := nanoTime(gameTime.CompetitiveStats.AllHeroes.Game.TimePlayed)
	if err != nil {
		logrus.Errorf("Error splitting time: %v", err)
		return nil, err
	}

	totalTime := quickTime + compTime
	return &totalTime, nil
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
		absTime = time.Duration(min) * time.Minute
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
		absTime = time.Duration(hour) * time.Hour
		min, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		absTime = time.Duration(min) * time.Minute
		sec, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}
		absTime = time.Duration(sec) * time.Second
	default:
		// returns error if parts of string exceeds 3 (hr:min:sec)
		logrus.Errorf("OW API has changed the way time is encoded")
		return 0, fmt.Errorf("OW API changed the way time is encoded, got |%d| parts", len(parts))
	}

	return absTime, nil
}