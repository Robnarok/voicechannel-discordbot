package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"voicebot-discord/m/config"
	"voicebot-discord/m/eventhandler"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config.ReadConfig()
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	eventhandler.Init()
	dg.AddHandler(eventhandler.VoiceChannelCreate)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
