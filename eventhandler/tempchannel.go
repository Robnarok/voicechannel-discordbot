package eventhandler

import (
	"fmt"
	"log"
	"voicebot-discord/m/config"

	"github.com/bwmarrin/discordgo"
)

var (
	m map[string]int
)

func Init() {
	m = make(map[string]int)
}

func checkChannelID(v *discordgo.VoiceStateUpdate) error {
	if v.ChannelID == "foobar" {
		return nil
	}
	return nil
}

func VoiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	logged := false
	user, _ := s.User(v.UserID)
	err := checkChannelID(v)
	if err != nil {
		v.ChannelID = ""
		log.Printf("Interessant\n")
	}
	if v.BeforeUpdate.ChannelID == v.ChannelID {
		return
	}

	//Debbug Zeugs
	affectedChannel, _ := s.Channel(v.ChannelID)
	affectedBeforChannel := ""

	if v.BeforeUpdate != nil {
		tmp, _ := s.Channel(v.BeforeUpdate.ChannelID)
		affectedBeforChannel = tmp.Name
	}

	if !logged && v.ChannelID == "" { //Disconnect
		log.Printf("%s disconnected aus %s\n", user.Username, affectedBeforChannel)
		logged = true
	}

	if !logged && v.BeforeUpdate == nil { //Initial Connect
		log.Printf("%s connectet in %s\n", user.Username, affectedChannel.Name)
		logged = true
	}

	if !logged && v.BeforeUpdate != nil && v.ChannelID != v.BeforeUpdate.ChannelID { //Channel gewechselt
		log.Printf("%s wechselt von %s zu %s\n", user.Username, affectedBeforChannel, affectedChannel.Name)
		logged = true
	}

	// Ende Debug Zeugs

	if v.ChannelID == config.Masterchannel {

		targetchannel, err := s.GuildChannelCreate(v.GuildID, user.Username, discordgo.ChannelTypeGuildVoice)
		m[targetchannel.ID] = 0
		log.Printf("Channel erstellt!\n")
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
			log.Printf("%s - %d\n", key, m[key])
			return
		}
	}

	if v.BeforeUpdate == nil {
		return
	}

	for key := range m {
		if v.BeforeUpdate.ChannelID == key {
			m[key] = m[key] - 1
			log.Printf("%s - %d\n", key, m[key])
			if m[key] == 0 {
				s.ChannelDelete(key)
				delete(m, key)
				log.Printf("%s destroyed\n", key)
				return
			}

		}
	}

}
