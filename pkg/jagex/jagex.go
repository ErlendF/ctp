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

// the varius types of runescape accounts
const (
	normal  = "normal"
	ironman = "ironman"
	hcim    = "hardcore ironman"
	uim     = "ultimate ironman"
)

// relevant urls to check the beforementiond types
var urls = map[string]string{
	normal:  "http://services.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s",
	ironman: "http://services.runescape.com/m=hiscore_oldschool_ironman/index_lite.ws?player=%s",
	hcim:    "http://services.runescape.com/m=hiscore_oldschool_hardcore_ironman/index_lite.ws?player=%s",
	uim:     "http://services.runescape.com/m=hiscore_oldschool_ultimate/index_lite.ws?player=%s",
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
	for i := 0; i < 23; i++ {
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

// validator for runescape username
func (j *Jagex) ValidateRSAccount(rsAcc *models.RunescapeAccount) error {
	matched, err := regexp.MatchString("^[A-Za-z0-9_ -]{1,12}$", rsAcc.Username)
	if err != nil {
		return err
	}
	if !matched {
		return models.NewReqErrStr("invalid Runescape account name", "invalid Runescape account name")
	}

	// unless otherwise specified, the account is assumed to be "normal"
	if rsAcc.AccountType == "" {
		rsAcc.AccountType = normal
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

	err = models.AccValStatusCode(resp.StatusCode, "Jagex", "invalid Runescape account name")
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseString := string(responseData)
	skills := strings.Split(responseString, "\n")
	totalFields := strings.Split(skills[0], ",")

	rsAcc.TotalLevel, err = strconv.Atoi(totalFields[1])
	if err != nil {
		return err
	}
	rsAcc.TotalXP, err = strconv.Atoi(totalFields[2])
	if err != nil {
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
	case normal:
		return xp / normalXPRates[i]
	case ironman:
		return xp / ironmanXPRates[i]
	case hcim:
		return xp / ironmanXPRates[i] // identical to ironman
	case uim:
		return xp / ultimateXPRates[i]
	}

	return 0
}
