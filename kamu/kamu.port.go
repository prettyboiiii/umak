package kamu

type Kamu interface {
	StartConversation() error
	GetPlaceInQueue(diaryNumber string, retryCount uint8) error
}
