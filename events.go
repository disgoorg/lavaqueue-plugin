package lavaqueue

import (
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

const (
	EventTypeQueueEnd lavalink.EventType = "QueueEndEvent"
)

type LavaQueueEventListener interface {
	OnQueueEnd(player disgolink.Player, event QueueEndEvent)
}

type QueueEndEvent struct {
	GuildID_ snowflake.ID `json:"guildId"`
}

func (QueueEndEvent) Op() lavalink.Op {
	return lavalink.OpEvent
}

func (QueueEndEvent) Type() lavalink.EventType {
	return EventTypeQueueEnd
}

func (e QueueEndEvent) GuildID() snowflake.ID {
	return e.GuildID_
}
