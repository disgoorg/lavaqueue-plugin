package lavaqueue

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

func GetQueue(ctx context.Context, node disgolink.Node, guildID snowflake.ID) (*Queue, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/v4/sessions/%s/players/%s/queue", node.SessionID(), guildID), nil)
	if err != nil {
		return nil, err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var queue Queue
	if err = unmarshalBody(rs, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

func UpdateQueue(ctx context.Context, node disgolink.Node, guildID snowflake.ID, queue QueueUpdate) (*lavalink.Track, error) {
	rqBody, err := marshalBody(queue)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("/v4/sessions/%s/players/%s/queue", node.SessionID(), guildID), rqBody)
	if err != nil {
		return nil, err
	}
	rq.Header.Add("Content-Type", "application/json")

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var track lavalink.Track
	if err = unmarshalBody(rs, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func QueueNextTrack(ctx context.Context, node disgolink.Node, guildID snowflake.ID, count int) (*lavalink.Track, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/next?count=%d", node.SessionID(), guildID, count), nil)
	if err != nil {
		return nil, err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var track lavalink.Track
	if err = unmarshalBody(rs, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func QueuePreviousTrack(ctx context.Context, node disgolink.Node, guildID snowflake.ID, count int) (*lavalink.Track, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/previous?count=%d", node.SessionID(), guildID, count), nil)
	if err != nil {
		return nil, err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var track lavalink.Track
	if err = unmarshalBody(rs, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func AddQueueTracks(ctx context.Context, node disgolink.Node, guildID snowflake.ID, tracks []QueueTrack) (*lavalink.Track, error) {
	rqBody, err := marshalBody(tracks)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/tracks", node.SessionID(), guildID), rqBody)
	if err != nil {
		return nil, err
	}
	rq.Header.Add("Content-Type", "application/json")

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	if rs.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var track lavalink.Track
	if err = unmarshalBody(rs, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func RemoveQueueTrack(ctx context.Context, node disgolink.Node, guildID snowflake.ID, trackID int) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/tracks/%d", node.SessionID(), guildID, trackID), nil)
	if err != nil {
		return err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)
}

func ShuffleQueue(ctx context.Context, node disgolink.Node, guildID snowflake.ID) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/shuffle", node.SessionID(), guildID), nil)
	if err != nil {
		return err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)
}

func ClearQueue(ctx context.Context, node disgolink.Node, guildID snowflake.ID) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/v4/sessions/%s/players/%s/queue", node.SessionID(), guildID), nil)
	if err != nil {
		return err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)

}

func GetHistory(ctx context.Context, node disgolink.Node, guildID snowflake.ID) ([]lavalink.Track, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/v4/sessions/%s/players/%s/history", node.SessionID(), guildID), nil)
	if err != nil {
		return nil, err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var history []lavalink.Track
	if err = unmarshalBody(rs, &history); err != nil {
		return nil, err
	}

	return history, nil
}

func ClearHistory(ctx context.Context, node disgolink.Node, guildID snowflake.ID) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/v4/sessions/%s/players/%s/history", node.SessionID(), guildID), nil)
	if err != nil {
		return err
	}

	rs, err := node.Rest().Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)
}

func marshalBody(v any) (io.Reader, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func unmarshalBody(rs *http.Response, v any) error {
	if rs.StatusCode < 200 || rs.StatusCode >= 300 {
		var lavalinkError lavalink.Error
		if err := json.NewDecoder(rs.Body).Decode(&lavalinkError); err != nil {
			return err
		}
		return lavalinkError
	}

	if rs.StatusCode == http.StatusNoContent {
		return nil
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(rs.Body).Decode(v)
}
