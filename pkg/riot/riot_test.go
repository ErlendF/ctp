package riot

import "net/http"

type mockClient struct {
	err error
}

func (m *mockClient) Do(req *http.Request) (resp *http.Response, err error) {



	return nil, nil
}