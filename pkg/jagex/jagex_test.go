package jagex

import (
	"ctp/pkg/models"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGetter struct {
	err error
}

func (m *mockGetter) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}
	resp.Body = ioutil.NopCloser(strings.NewReader(testData))
	return resp, nil
}

var testData = `6355,2277,387381708
57292,99,15553384
20982,99,18653311
73125,99,15691909
48610,99,26081909
54148,99,21211037
30199,99,13041577
10689,99,22485307
2342,99,46913553
17900,99,13760137
81027,99,13034656
27850,99,13050488
57625,99,13053422
21428,99,13164200
3233,99,14953077
9975,99,13551444
13733,99,13101666
7391,99,13222039
27716,99,13088572
29789,99,13709440
45744,99,13142445
826,99,20648671
19558,99,13039828
3494,99,13229636
-1,-1
-1,-1
-1,-1
281553,51
-1,-1
-1,-1
-1,-1
503745,1
192248,2
8499,48
-1,-1`

func TestGetRSPlaytime(t *testing.T) {
	var cases = []struct {
		name        string
		getterErr   error
		expectedErr error
	}{
		{"Test ok", nil, nil},
		{"Test getter error", errors.New("test"), errors.New("test")},
	}

	mg := &mockGetter{}
	jagex := New(mg)

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Initializing mock structs with random data
			var rsAcc models.RunescapeAccount
			err := faker.FakeData(&rsAcc.Username)
			require.Nil(t, err)
			rsAcc.AccountType = "normal"

			mg.err = tc.getterErr

			game, err := jagex.GetRSPlaytime(&rsAcc)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				if assert.NotNil(t, game) {
					assert.Equal(t, "Old School Runescape", game.Name)
				}
			}
		})
	}
}

func TestValidateRSAccount(t *testing.T) {
	var cases = []struct {
		name        string
		rsAcc       models.RunescapeAccount
		getterErr   error
		expectedErr error
	}{
		{"Test ok", models.RunescapeAccount{Username: "Test123", AccountType: "normal"}, nil, nil},
		{"Test ok ironman", models.RunescapeAccount{Username: "Test123", AccountType: "ironman"}, nil, nil},
		{"Test ok hardcore ironman", models.RunescapeAccount{Username: "Test123", AccountType: "hardcore ironman"}, nil, nil},
		{"Test ok ultimate ironman", models.RunescapeAccount{Username: "Test123", AccountType: "ultimate ironman"}, nil, nil},
		{"Test invalid account type", models.RunescapeAccount{Username: "Test123", AccountType: "this is invalid"}, nil,
			models.NewReqErrStr("invalid Runescape account type", "invalid Runescape account type")},
		{"Test getter error", models.RunescapeAccount{Username: "Test123", AccountType: "normal"}, errors.New("test"), errors.New("test")},
		{"Test too long name error", models.RunescapeAccount{Username: "this name is too long", AccountType: "normal"}, nil,
			models.NewReqErrStr("invalid Runescape account name", "invalid Runescape account name")},
	}

	mg := &mockGetter{}
	jagex := New(mg)

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mg.err = tc.getterErr

			err := jagex.ValidateRSAccount(&tc.rsAcc)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
