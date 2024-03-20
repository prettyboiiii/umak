package kamu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	HOST              = "https://networkmigri.boost.ai"
	BASE_API_PATH     = HOST + "/api"
	CHAT_API_PATH     = "chat/v2"
	MY_PLACE_IN_Queue = "My place in queue"
)

var (
	SeesionEndedErr error = errors.New("conversation.ended")
	MAX_RETRY       uint8 = 3
)

type kamu struct {
	logger               log.Logger
	ConversationID       string
	PlaceInQueueChoiceId string
}

func init() {
	maxRetry, isFound := os.LookupEnv("MAX_RETRY")
	if isFound {
		if maxRetry, err := strconv.ParseUint(maxRetry, 10, 8); err == nil {
			MAX_RETRY = uint8(maxRetry)
		}
	}
}

// GetPlaceInQueue implements Kamu.
func (k *kamu) GetPlaceInQueue(diaryNumber string, retryCount uint8) error {
	if k.ConversationID == "" || k.PlaceInQueueChoiceId == "" {
		if err := k.StartConversation(); err != nil {
			return err
		}
	}
	data := &request{
		Command:        commandPost,
		ConversationID: k.ConversationID,
		ID:             k.PlaceInQueueChoiceId,
		Type:           requestTypeActionLink,
	}
	resp, err := newRequest[conversationResponse](CHAT_API_PATH, http.MethodPost, data)
	if err != nil {
		return err
	}

	for _, element := range resp.Response.Elements {
		if strings.Contains(element.Payload.HTML, "Enter your diary number") {
			data.ID = ""
			data.Value = diaryNumber
			data.Type = requestTypeText
			k.logger.Printf("...Checking Diary number: %s ...\n", diaryNumber)
			resp, err := newRequest[conversationResponse](CHAT_API_PATH, http.MethodPost, data)
			if err != nil {
				return err
			}
			for _, element := range resp.Response.Elements {
				if element.Payload.JSON != nil {
					k.logger.Println("The application status with this diary number is:", element.Payload.JSON.Data.CounterValue)
					return nil
				}
				if retryCount <= MAX_RETRY && strings.Contains(element.Payload.HTML, "Enter your diary number") {
					return k.GetPlaceInQueue(diaryNumber, retryCount+1)
				}
			}
		}
	}

	return nil
}

// StartConversation implements Kamu.
func (k *kamu) StartConversation() error {
	var data = &request{
		Command:          commandStart,
		DisableHumanChat: false,
	}
	res, err := newRequest[conversationResponse](CHAT_API_PATH, http.MethodPost, data)
	if err != nil {
		return err
	}
	k.ConversationID = res.Conversation.ID
	if k.ConversationID == "" {
		return errors.New("Cannot acquire ConversationID")
	}

	for _, element := range res.Response.Elements {
		if len(element.Payload.Links) == 0 {
			continue
		}
		for _, link := range element.Payload.Links {
			if link.Text == MY_PLACE_IN_Queue {
				k.PlaceInQueueChoiceId = link.ID
				return nil
			}
		}
	}
	if k.PlaceInQueueChoiceId == "" {
		return errors.New("Cannot acquire PlaceInQueueChoiceId")
	}

	return nil
}

func New() Kamu {
	return &kamu{}
}

func newRequest[R any](
	path string,
	method string,
	data *request,
) (*R, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", BASE_API_PATH, path), data.NewReader())
	if err != nil {
		log.Fatal(err)
	}

	setRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	errorRes := new(errorResponse)
	if err := json.Unmarshal(bodyText, errorRes); err == nil && *errorRes != (errorResponse{}) {
		return nil, errors.New(errorRes.Tag)
	}

	result := new(R)
	if err := json.Unmarshal(bodyText, result); err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func setRequestHeaders(outR *http.Request) {
	outR.Header.Set("Accept", "application/json, text/plain, */*")
	outR.Header.Set("Accept-Language", "en-US,en;q=0.7,th;q=0.3")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	outR.Header.Set("Content-Type", "application/json")
	outR.Header.Set("Origin", HOST)
	outR.Header.Set("DNT", "1")
	outR.Header.Set("Sec-GPC", "1")
	outR.Header.Set("Connection", "keep-alive")
	outR.Header.Set("Referer", HOST)
	outR.Header.Set("Sec-Fetch-Dest", "empty")
	outR.Header.Set("Sec-Fetch-Mode", "cors")
	outR.Header.Set("Sec-Fetch-Site", "same-origin")
	outR.Header.Set("TE", "trailers")
}
