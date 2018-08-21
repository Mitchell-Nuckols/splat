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
	Text        string `json:"text,omitempty"`
	Attachments `json:"attachments,omitempty"`
	Markdown    bool `json:"mrkdwn,omitempty"`
}

// Attachments are part of extra data that can be sent in a Slack response
type Attachments []struct {
	Fallback   string `json:"fallback,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Color      string `json:"color,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Pretext    string `json:"pretext,omitempty"`
	Text       string `json:"text,omitempty"`
	Fields     `json:"fields,omitempty"`
	ImageURL   string   `json:"image_url,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
	Timestamp  int      `json:"ts,omitempty"`
}

// Fields are a part of Attachments
type Fields []struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}
