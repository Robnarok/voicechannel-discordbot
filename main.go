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

var (
	m map[string]int
)

func main() {
	config.ReadConfig()
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	m = make(map[string]int)

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

	if v.ChannelID == config.Masterchannel {
		user, _ := s.User(v.UserID)

		targetchannel, err := s.GuildChannelCreate(v.GuildID, user.Username, discordgo.ChannelTypeGuildVoice)
		m[targetchannel.ID] = 0
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//s.ChannelVoiceJoin(v.GuildID, targetchannel.ID, false, false)
		err = s.GuildMemberMove(v.GuildID, v.UserID, &targetchannel.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		time.Sleep(1 * time.Second)

	}

	for key := range m {
		if v.ChannelID == key {
			m[key] = m[key] + 1
		}
	}

	if v.BeforeUpdate != nil {
		for key := range m {
			if v.BeforeUpdate.ChannelID != v.ChannelID {
				if v.BeforeUpdate.ChannelID == key {
					m[key] = m[key] - 1
					if m[key] == 0 {
						s.ChannelDelete(key)
						delete(m, key)
						return
					}
				}
			}
		}
		fmt.Println(m)
	}

}
