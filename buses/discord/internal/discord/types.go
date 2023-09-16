package discord

// incomplete, obviously
type Message struct {
	ID        string
	Author    User
	Content   string
	ChannelID string `mapstructure:"channel_id"`
}

type User struct {
	ID            string
	Username      string
	Discriminator string
	IsBot         bool `mapstructure:"bot"`
}
