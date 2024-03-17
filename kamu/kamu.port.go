package kamu

type Kamu interface {
	StartConversation() error
	GetPlaceInQueue(diaryNumber string) error
}
