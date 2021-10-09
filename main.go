package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"voicebot-discord/m/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config.ReadConfig()
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(voiceChannelCreate)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func voiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID == "835121335851155470" {
		foo, _ := s.User(v.UserID)
		targetchannel, err := s.GuildChannelCreate(v.GuildID, foo.Username, discordgo.ChannelTypeGuildVoice)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("channelid: %s \n userid: %s \n", v.ChannelID, v.UserID)
		time.Sleep(250 * time.Millisecond)
		s.ChannelVoiceJoin(v.GuildID, targetchannel.ID, false, false)
		err = s.GuildMemberMove(v.GuildID, v.UserID, &targetchannel.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		time.Sleep(10 * time.Second)
		s.ChannelDelete(targetchannel.ID)

	}
}
