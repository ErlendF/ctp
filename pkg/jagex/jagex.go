package jagex

import (
	"ctp/pkg/models"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Jagex is a struct which contains everything necessary to handle a request related to Jagex
type Jagex struct {
	models.Getter
}

// New returns a new Jagex instance
func New(getter models.Getter) *Jagex {
	return &Jagex{getter}
}

const normalHiscores = "http://services.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s"

var urls = map[string]string{
	"normal":           "http://services.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s",
	"ironman":          "http://services.runescape.com/m=hiscore_oldschool_ironman/index_lite.ws?player=%s",
	"hardcore ironman": "http://services.runescape.com/m=hiscore_oldschool_hardcore_ironman/index_lite.ws?player=%s",
	"ultimate ironman": "http://services.runescape.com/m=hiscore_oldschool_ultimate/index_lite.ws?player=%s",
}

// GetRSPlaytime returns an estimate for time spent playing Runescape
func (j *Jagex) GetRSPlaytime(rsAcc *models.RunescapeAccount) (*models.Game, error) {
	url, ok := urls[rsAcc.AccountType]

	if !ok {
		return nil, fmt.Errorf("invalid account type in GetRSPlaytime: %s", rsAcc.AccountType)
	}

	response, err := j.Get(fmt.Sprintf(url, rsAcc.Username))
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
	lines := strings.Split(responseString, "\n")
	for i := 1; i < 23; i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) != 3 {
			return nil, errors.New("wrong number of fields in GetRSPlaytime")
		}

		xp, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}

		time += xpToTime(xp, i, rsAcc.AccountType)
	}

	return &models.Game{Time: time, Name: "Runescape"}, nil
}

func (j *Jagex) ValidateRSAccount(rsAcc *models.RunescapeAccount) error {
	matched, err := regexp.MatchString("^[A-Za-z0-9_ -]{1,12}$", rsAcc.Username)
	if err != nil {
		return err
	}

	if !matched {
		return models.NewReqErrStr("invalid Runescape account name", "invalid Runescape account name")
	}

	if rsAcc.AccountType == "" {
		rsAcc.AccountType = "normal"
	}

	url, ok := urls[rsAcc.AccountType]

	if !ok {
		return models.NewReqErrStr("invalid Runescape account type", "invalid Runescape account type")
	}

	resp, err := j.Get(fmt.Sprintf(url, rsAcc.Username))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := models.AccValStatusCode(resp.StatusCode, "Jagex", "invalid Runescape account name"); err != nil {
		return err
	}

	return nil
}

// xpToTime estimates time spent on one skill based on the xp (Experience Points)
func xpToTime(xp, i int, accountType string) int {
	if i < 1 || i > 23 {
		return 0
	}

	switch accountType {
	case "normal":
		return xp / normalXPRates[i]
	case "ironman":
		return xp / ironmanXPRates[i]
	case "hardcore ironman":
		return xp / ironmanXPRates[i] // identical to ironman
	case "ultimate ironman":
		return xp / ultimateXPRates[i]
	}

	return 0
}

// normalXPRates for each skill
var normalXPRates = [...]int{
	0,
	90000,
	90000,
	90000,
	300000,
	150000,
	200000,
	100000,
	400000,
	70000,
	250000,
	70000,
	200000,
	150000,
	250000,
	60000,
	200000,
	44000,
	100000,
	50000,
	100000,
	50000,
	120000,
	400000,
}

// ironmanXPRates for each skill
var ironmanXPRates = [...]int{
	0,
	90000,
	90000,
	90000,
	300000,
	150000,
	200000,
	100000,
	400000,
	70000,
	250000,
	70000,
	200000,
	150000,
	250000,
	60000,
	200000,
	44000,
	100000,
	50000,
	100000,
	50000,
	120000,
	400000,
}

// ultimateXPRates for each skill
var ultimateXPRates = [...]int{
	0,
	90000,
	90000,
	90000,
	300000,
	150000,
	200000,
	100000,
	400000,
	70000,
	250000,
	70000,
	200000,
	150000,
	250000,
	60000,
	200000,
	44000,
	100000,
	50000,
	100000,
	50000,
	120000,
	400000,
}
