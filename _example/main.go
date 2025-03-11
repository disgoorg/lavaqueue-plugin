package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/lavaqueue-plugin"
)

var (
	token   = os.Getenv("BOT_TOKEN")
	guildID = snowflake.GetEnv("GUILD_ID")

	nodeName     = os.Getenv("LAVALINK_NODE_NAME")
	nodeAddress  = os.Getenv("LAVALINK_NODE_ADDRESS")
	nodePassword = os.Getenv("LAVALINK_NODE_PASSWORD")
	nodeSecure   = os.Getenv("LAVALINK_NODE_SECURE")
)

type exampleBot struct {
	Discord  bot.Client
	Lavalink disgolink.Client
}

func main() {
	slog.Info("starting example bot...")

	b := &exampleBot{}

	r := handler.New()
	r.Use(middleware.Go)
	r.SlashCommand("/play", b.onPlay)

	// create a new bot
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListeners(r),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		slog.Error("error while creating bot", slog.Any("err", err))
		return
	}
	b.Discord = client

	b.Lavalink = disgolink.New(client.ApplicationID(),
		disgolink.WithPlugins(lavaqueue.New()),
		disgolink.WithListenerFunc(b.onQueueEnd),
	)

	var guildIDs []snowflake.ID
	if guildID > 0 {
		guildIDs = append(guildIDs, guildID)
	}
	if err = handler.SyncCommands(client, commands, guildIDs); err != nil {
		slog.Error("error while syncing commands", slog.Any("err", err))
	}

	disgoCtx, disgoCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer disgoCancel()
	if err = client.OpenGateway(disgoCtx); err != nil {
		slog.Error("error while opening gateway", slog.Any("err", err))
		return
	}
	defer client.Close(context.Background())

	lavalinkCtx, lavalinkCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer lavalinkCancel()
	if _, err = b.Lavalink.AddNode(lavalinkCtx, disgolink.NodeConfig{
		Name:     nodeName,
		Address:  nodeAddress,
		Password: nodePassword,
		Secure:   nodeSecure == "true",
	}); err != nil {
		slog.Error("error while adding node", slog.Any("err", err))
	}

	slog.Info("example bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
