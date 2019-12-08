package gorest

type Session struct {
	Secret string `json:"jwt_secret"`
	UserID string
}

