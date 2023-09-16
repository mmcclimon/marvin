package discord

type OpType int

const (
	// An event was dispatched
	Dispatch OpType = iota
	// Fired periodically by the client to keep the connection alive.
	Heartbeat
	// Starts a new session during the initial handshake.
	Identify
	// Update the client's presence.
	PresenceUpdate
	// Used to join/leave or move between voice channels.
	VoiceStateUpdate
	//nolint:unused
	unused // 5 is just missing
	// Resume a previous session that was disconnected.
	Resume
	// You should attempt to reconnect and resume immediately.
	Reconnect
	// Request information about offline guild members in a large guild.
	RequestGuildMembers
	// The session has been invalidated. You should reconnect and identify/resume accordingly.
	InvalidSession
	// Sent immediately after connecting, contains the heartbeat_interval to use.
	Hello
	// Sent in response to receiving a heartbeat to acknowledge that it has been received.
	HeartbeatACK
)

type GatewayEvent struct {
	Op   OpType
	Data any       `json:"d"`
	Seq  *int      `json:"s"`
	Type EventType `json:"t,omitempty"`
}

type EventType string

const (
	MessageCreate EventType = "MESSAGE_CREATE"
)
