package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

func voiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID == config.Masterchannel {
		user, _ := s.User(v.UserID)

		targetchannel, err := s.GuildChannelCreate(v.GuildID, user.Username, discordgo.ChannelTypeGuildVoice)
		m[targetchannel.ID] = 0
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := discordgo.ChannelEdit{
			ParentID: config.KategoriId,
			Position: 1,
		}
		s.ChannelEditComplex(targetchannel.ID, &data)

		err = s.GuildMemberMove(v.GuildID, v.UserID, &targetchannel.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	for key := range m {
		if v.ChannelID == key {
			m[key] = m[key] + 1
		}
	}

	if v.BeforeUpdate == nil {
		return
	}

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

}
