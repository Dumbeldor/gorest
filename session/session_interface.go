package gorest

type ReaderInterface interface {
	LoadSession(userID string) *Session
}
