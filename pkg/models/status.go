package models

//StatusResp contains everything returned by the /status path
type StatusResp struct {
	Gitlab  int     `json:"gitlab"`
	DB      int     `json:"database"`
	Uptime  float64 `json:"uptime"`
	Version string  `json:"version"`
}

//Status is an interface which defines all methods a "Status" should provide
type Status interface {
	GetStatus() (*StatusResp, error)
}
