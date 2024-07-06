package lavaqueue

import (
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
)

var (
	_ disgolink.EventPlugins = (*Plugin)(nil)
	_ disgolink.Plugin       = (*Plugin)(nil)
)

func New() *Plugin {
	return NewWithLogger(slog.Default())
}

func NewWithLogger(logger *slog.Logger) *Plugin {
	return &Plugin{
		eventPlugins: []disgolink.EventPlugin{
			&queueEndHandler{
				logger: logger,
			},
		},
	}
}

type Plugin struct {
	eventPlugins []disgolink.EventPlugin
}

func (p *Plugin) EventPlugins() []disgolink.EventPlugin {
	return p.eventPlugins
}

func (p *Plugin) Name() string {
	return "lavaqueue"
}

func (p *Plugin) Version() string {
	return "1.0.0"
}

var _ disgolink.EventPlugin = (*queueEndHandler)(nil)

type queueEndHandler struct {
	logger *slog.Logger
}

func (h *queueEndHandler) Event() lavalink.EventType {
	return EventTypeQueueEnd
}

func (h *queueEndHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e QueueEndEvent
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Error("Failed to unmarshal QueueEndEvent", slog.Any("err", err))
		return
	}

	player.Lavalink().EmitEvent(player, e)
}
