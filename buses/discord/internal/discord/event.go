package discord

import "nhooyr.io/websocket"

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
	TypeMessageCreate EventType = "MESSAGE_CREATE"
	TypeReady         EventType = "READY"
	TypeResumed       EventType = "RESUMED"
)

// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
type CloseCode int

const (
	UnknownError         websocket.StatusCode = 4000
	UnknownOpcode        websocket.StatusCode = 4001
	DecodeError          websocket.StatusCode = 4002
	NotAuthenticated     websocket.StatusCode = 4003
	AuthenticationFailed websocket.StatusCode = 4004
	AlreadyAuthenticated websocket.StatusCode = 4005
	InvalidSeq           websocket.StatusCode = 4007
	RateLimited          websocket.StatusCode = 4008
	SessionTimedOut      websocket.StatusCode = 4009
	InvalidShard         websocket.StatusCode = 4010
	ShardingRequired     websocket.StatusCode = 4011
	InvalidAPIVersion    websocket.StatusCode = 4012
	InvalidIntent        websocket.StatusCode = 4013
	DisallowedIntent     websocket.StatusCode = 4014
)
