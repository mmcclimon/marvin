package discord

// incomplete, obviously
type Message struct {
	ID      string
	Author  User
	Content string
}

type User struct {
	ID            string
	Username      string
	Discriminator string
	IsBot         bool `mapstructure:"bot"`
}
