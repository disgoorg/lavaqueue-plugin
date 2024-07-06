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

func GetQueue(ctx context.Context, client disgolink.RestClient, sessionID string, guildID snowflake.ID) (*Queue, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/v4/sessions/%s/players/%s/queue", sessionID, guildID), nil)
	if err != nil {
		return nil, err
	}

	rs, err := client.Do(rq)
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

func UpdateQueue(ctx context.Context, client disgolink.RestClient, sessionID string, guildID snowflake.ID, queue QueueUpdate) error {
	rqBody, err := marshalBody(queue)
	if err != nil {
		return err
	}

	rq, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("/v4/sessions/%s/players/%s/queue", sessionID, guildID), rqBody)
	if err != nil {
		return err
	}

	rs, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)
}

func AddQueueTracks(ctx context.Context, client disgolink.RestClient, sessionID string, guildID snowflake.ID, tracks []QueueTrack) error {
	rqBody, err := marshalBody(tracks)
	if err != nil {
		return err
	}

	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/v4/sessions/%s/players/%s/queue/tracks", sessionID, guildID), rqBody)
	if err != nil {
		return err
	}

	rs, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	return unmarshalBody(rs, nil)
}

func GetHistory(ctx context.Context, client disgolink.RestClient, sessionID string, guildID snowflake.ID) ([]lavalink.Track, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/v4/sessions/%s/players/%s/history", sessionID, guildID), nil)
	if err != nil {
		return nil, err
	}

	rs, err := client.Do(rq)
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

	return json.NewDecoder(rs.Body).Decode(v)
}
