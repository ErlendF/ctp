package models

//UserManager contains all functions a usermanager is expected to provide
type UserManager interface {
	GetUserInfo(username string) (*UserInfo, error)
}
