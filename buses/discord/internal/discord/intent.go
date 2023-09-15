package discord

type Intent int

// taken from https://discord.com/developers/docs/topics/gateway#list-of-intents
const (
	Guilds Intent = 1 << iota
	GuildMembers
	GuildModeration
	GuildEmojiAndStickers
	GuildIntegrations
	GuildWebhooks
	GuildInvites
	GuildVoiceStates
	GuildPresences
	GuildMessages
	GuildMessageReactions
	GuildMessageTyping
	DirectMessages
	DirectMessageReactions
	DirectMessageTyping
	MessageContent
	GuildScheduledEvents
	AutoModerationConfiguration = 1 << 20
	AutoModerationExecution     = 1 << 21
)
