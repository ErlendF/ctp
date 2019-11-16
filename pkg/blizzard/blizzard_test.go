package blizzard

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockBlizzard struct {
	setup respSetup
}

func (m *mockBlizzard) Get(url string) (*http.Response, error){
	validUrl := strings.Contains(url, "ow-api.com/v1/stats/")
	if !validUrl {
		return nil, m.setup.err
	}
	resp := &http.Response{StatusCode:http.StatusOK, Header:make(http.Header, 0)}
	body, err := json.Marshal(m.setup.resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return resp, nil
}

type respSetup struct {
	/*CompetitiveStats struct {
		CareerStats struct {
			AllHeroes struct {
				Game struct {
					TimePlayed string
				}
			}
		}
	}

	*/
	/*QuickPlayStats   struct {
		CareerStats struct {
			AllHeroes struct {
				Game struct {
					TimePlayed string
				}
			}
		}
	}

	 */
	resp             models.BlizzardResp
	status 			 int
	err              error
}

func TestBlizzard_GetBlizzardPlaytime(t *testing.T) {
	var testcase = []struct{
		name          string
		payload       *models.Overwatch
		cTime         string
		qTime         string
		expectedError error
	}{
		{"Test OK",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "", Region:    ""}, "","",fmt.Errorf("")},
		{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "", Region:    ""}, "","",fmt.Errorf("")},
		{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "", Region:    ""}, "","",fmt.Errorf("")},
		{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "", Region:    ""}, "","",fmt.Errorf("")},
	}

	//client := &mockBlizzard{}
	//ow := New(client)

	for _, item := range testcase {
		t.Run(item.name, func(t *testing.T){

			setup := respSetup{}
			setup.resp.CompetitiveStats.CareerStats.AllHeroes.Game.TimePlayed = item.cTime
			setup.resp.QuickPlayStats.CareerStats.AllHeroes.Game.TimePlayed = item.qTime



		})
	}

}