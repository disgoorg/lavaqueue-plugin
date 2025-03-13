package main

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"

	"github.com/disgoorg/lavaqueue-plugin"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "play",
		Description: "Plays a song",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "identifier",
				Description: "The song/query to play",
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "next",
		Description: "Skips the current song",
	},
	discord.SlashCommandCreate{
		Name:        "last",
		Description: "Plays the last played song",
	},
	discord.SlashCommandCreate{
		Name:        "queue",
		Description: "Shows the current queue",
	},
	discord.SlashCommandCreate{
		Name:        "pause",
		Description: "Pauses the current song",
	},
	discord.SlashCommandCreate{
		Name:        "resume",
		Description: "Resumes the current song",
	},
	discord.SlashCommandCreate{
		Name:        "stop",
		Description: "Stops the current song",
	},
}

func (b *exampleBot) onPlay(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	identifier := data.String("identifier")
	if source, ok := data.OptString("source"); ok {
		identifier = lavalink.SearchType(source).Apply(identifier)
	} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	voiceState, ok := b.Discord.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command",
		})
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := b.Lavalink.BestNode().LoadTracks(ctx, identifier)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("error while loading tracks"),
		})
		return err
	}

	var tracks []lavalink.Track
	switch result.LoadType {
	case lavalink.LoadTypeTrack:
		tracks = []lavalink.Track{result.Data.(lavalink.Track)}
	case lavalink.LoadTypePlaylist:
		tracks = result.Data.(lavalink.Playlist).Tracks
	case lavalink.LoadTypeSearch:
		tracks = result.Data.(lavalink.Search)
	case lavalink.LoadTypeEmpty:
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("no tracks found"),
		})
		return err
	case lavalink.LoadTypeError:
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("error while loading tracks"),
		})
		return err
	}

	queueTracks := make([]lavaqueue.QueueTrack, len(tracks))
	for i := range tracks {
		queueTracks[i] = lavaqueue.QueueTrack{
			Encoded:  tracks[i].Encoded,
			UserData: tracks[i].UserData,
		}
	}

	if err = b.Discord.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, false); err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("error while connecting to voice channel"),
		})
		return err
	}

	player := b.Lavalink.Player(*e.GuildID())
	nextTrack, err := lavaqueue.AddQueueTracks(context.TODO(), player.Node(), *e.GuildID(), queueTracks)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("error while adding tracks to queue"),
		})
		return err
	}

	if nextTrack != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Playing: %s", nextTrack.Info.Title)),
		})
		return err
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: json.Ptr("Added to queue"),
	})
	return nil
}

func (b *exampleBot) onNext(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	player := b.Lavalink.Player(*e.GuildID())
	if player == nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Not connected to a voice channel",
		})
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	track, err := lavaqueue.QueueNextTrack(e.Ctx, player.Node(), *e.GuildID())
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("error while skipping track"),
		})
		return err
	}

	if track == nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("no tracks in queue"),
		})
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: json.Ptr(fmt.Sprintf("Playing: %s", track.Info.Title)),
	})
	return err
}
