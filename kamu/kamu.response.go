package kamu

type jsonPayload struct {
	Template string `json:"template"`
	Data     struct {
		CounterValue  string `json:"counterValue"`
		MarkerColor   string `json:"markerColor"`
		LineThickness string `json:"lineThickness"`
		LinePosition  string `json:"linePosition"`
		Messages      struct {
			StartText   string `json:"start.text"`
			CounterText string `json:"counter.text"`
			EndText     string `json:"end.text"`
		} `json:"messages"`
		LineColor       string `json:"lineColor"`
		MarkerThickness string `json:"markerThickness"`
	} `json:"data"`
}

type linkPayload struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type payload struct {
	HTML  string        `json:"html"`
	Links []linkPayload `json:"links"`
	JSON  *jsonPayload  `json:"json"`
}

type conversationResponse struct {
	Response struct {
		Source   string `json:"source"`
		Elements []struct {
			Type    string  `json:"type"`
			Payload payload `json:"payload,omitempty"`
		} `json:"elements"`
		Language    string `json:"language"`
		AvatarURL   string `json:"avatar_url"`
		DateCreated string `json:"date_created"`
		ID          string `json:"id"`
	} `json:"response"`
	Conversation struct {
		ID        string `json:"id"`
		Reference string `json:"reference"`
		State     struct {
			AllowDeleteConversation bool   `json:"allow_delete_conversation"`
			ChatStatus              string `json:"chat_status"`
			MaxInputChars           int    `json:"max_input_chars"`
			PrivacyPolicyURL        string `json:"privacy_policy_url"`
		} `json:"state"`
	} `json:"conversation"`
}

type errorResponse struct {
	Error string `json:"error"`
	Type  string `json:"type"`
	Tag   string `json:"tag"`
}
