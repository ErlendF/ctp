package jagex

import (
	"ctp/pkg/models"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Jagex is a struct which contains everything necessary to handle a request related to valve
type Jagex struct {
	models.Getter
}

// New returns a new valve instance
func New(getter models.Getter) *Jagex {
	return &Jagex{getter}
}

const normalHiscores = "http://services.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s"

// GetRSPlaytime returns an estimate for time spent playing Runescape
func (j *Jagex) GetRSPlaytime(username string) (*models.Game, error) {
	response, err := j.Get(fmt.Sprintf(normalHiscores, username))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var time int
	responseString := string(responseData)

	lines := strings.Split(string(responseString), "\n")
	for i := 0; i < 24; i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) != 3 {
			return nil, errors.New("wrong number of fields in GetRSPlaytime")
		}

		xp := fields[2]
		estimate, err := xpToTime(xp, i)
		if err != nil {
			return nil, err
		}

		time += estimate
	}

	game := &models.Game{Time: time, Name: "Runescape"}

	logrus.Debugf("rs game: %+v", game)
	return game, nil
}

func xpToTime(xp string, index int) (int, error) {
	input, err := strconv.Atoi(xp)
	if err != nil {
		return 0, err
	}
	var output int
	switch index {
	case 1: //Attack
		output = input / 90000
		break
	case 2: //Defence
		output = input / 90000
		break
	case 3: //Strength
		output = input / 90000
		break
	case 4: //Hitpoints
		output = input / 300000
		break
	case 5: //Ranged
		output = input / 150000
		break
	case 6: ///Prayer
		output = input / 200000
		break
	case 7: //Magic
		output = input / 100000
		break
	case 8: //Cooking
		output = input / 400000
		break
	case 9: //Woodcutting
		output = input / 70000
		break
	case 10: //Fletching
		output = input / 250000
		break
	case 11: //Fishing
		output = input / 70000
		break
	case 12: //Firemaking
		output = input / 200000
		break
	case 13: //Crafting
		output = input / 150000
		break
	case 14: //Smithing
		output = input / 250000
		break
	case 15: //Mining
		output = input / 60000
		break
	case 16: //Herblore
		output = input / 200000
		break
	case 17: //Agility
		output = input / 44000
		break
	case 18: //Thieving
		output = input / 100000
		break
	case 19: //Slayer
		output = input / 50000
		break
	case 20: //Farming
		output = input / 100000
		break
	case 21: //Runecraft
		output = input / 50000
		break
	case 22: //Hunter
		output = input / 120000
		break
	case 23: //Construction
		output = input / 400000
		break
	}

	fmt.Println("Index: " + strconv.Itoa(index) + "\nOutput: " + xp)
	return output, nil
}
