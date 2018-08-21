package splat

// Command creates a new command to be executed
type Command struct {
	Name    string
	Execute func(*Payload) *Response
}

// Payload is the data recieved from Slack
type Payload struct {
	Token,
	TeamID,
	TeamDomain,
	EnterpriseID,
	EnterpriseName,
	ChannelID,
	ChannelName,
	UserID,
	UserName,
	Command,
	Text,
	ResponseURL,
	TriggerID string
}

// Response is the data sent back to Slack
type Response struct {
	Text string `json:"text"`
}
