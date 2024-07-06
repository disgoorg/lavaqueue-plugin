package lavaqueue

import (
	"github.com/disgoorg/disgolink/v3/lavalink"
)

type QueueType string

const (
	QueueTypeNormal      QueueType = "normal"
	QueueTypeRepeatTrack QueueType = "repeat_track"
	QueueTypeRepeatQueue QueueType = "repeat_queue"
)

type Queue struct {
	Type   QueueType        `json:"type"`
	Tracks []lavalink.Track `json:"tracks"`
}

type QueueTrack struct {
	Encoded  string           `json:"encoded"`
	UserData lavalink.RawData `json:"user_data"`
}

type QueueUpdate struct {
	Type   *QueueType    `json:"type,omitempty"`
	Tracks *[]QueueTrack `json:"tracks,omitempty"`
}
