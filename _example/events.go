package main

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"

	"github.com/disgoorg/lavaqueue-plugin"
)

func (b *exampleBot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Discord.ApplicationID() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	b.Lavalink.OnVoiceStateUpdate(ctx, event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
}

func (b *exampleBot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	b.Lavalink.OnVoiceServerUpdate(ctx, event.GuildID, event.Token, *event.Endpoint)
}

func (b *exampleBot) onQueueEnd(player disgolink.Player, event lavaqueue.QueueEndEvent) {
	// do something
}
