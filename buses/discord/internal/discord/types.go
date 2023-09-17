package discord

// incomplete, obviously
type Message struct {
	ID        string
	Author    User
	Content   string
	ChannelID string `mapstructure:"channel_id"`
	Mentions  []User
}

type User struct {
	ID            string
	Username      string
	Discriminator string
	IsBot         bool `mapstructure:"bot"`
}

type Ready struct {
	APIVersion       int `mapstructure:"v"`
	User             User
	SessionID        string `mapstructure:"session_id"`
	ResumeGatewayURL string `mapstructure:"resume_gateway_url"`

	// ignoring, for now: guilds, application
}

// payload for sending resume ops
type GatewayResume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       *int   `json:"seq"`
}
