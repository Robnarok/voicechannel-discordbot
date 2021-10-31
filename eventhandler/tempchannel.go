package eventhandler

import (
	"fmt"
	"log"
	"voicebot-discord/m/config"

	"github.com/bwmarrin/discordgo"
)

type GeneratedChannel struct {
	kategory     string
	voicechannel string
	temprole     string
	textchannel  string
	current      int
}

var (
	m map[string]GeneratedChannel
)

func Init() {
	m = make(map[string]GeneratedChannel)
}

func checkChannelID(v *discordgo.VoiceStateUpdate) error {
	if v.ChannelID == "foobar" {
		return nil
	}
	return nil
}

func manipulatePermissions(s *discordgo.Session, voicechannel string, textchannel string, kategorie string, user string, rolle string) {
	everybode := "835121335851155466"

	// Jeder in der ROlle darf den Textchannel sehen
	s.ChannelPermissionSet(
		textchannel,
		rolle,
		discordgo.PermissionOverwriteTypeRole,
		1024,
		0)
	// Jeder ohne Rolle darf den Textchannel nicht sehen
	s.ChannelPermissionSet(
		textchannel,
		everybode,
		discordgo.PermissionOverwriteTypeRole,
		0,
		1024)
	// Ersteller bekommt Admin rechte auf die Kategorie
	s.ChannelPermissionSet(
		kategorie,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)
	// Ersteller auf Textchannel
	s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)
	// Ersteller auf Voicechannel
	s.ChannelPermissionSet(
		voicechannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)

}

func VoiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	logged := false
	user, _ := s.User(v.UserID)
	err := checkChannelID(v)
	if err != nil {
		v.ChannelID = ""
		log.Printf("Interessant\n")
	}

	//Debbug Zeugs
	affectedChannel, _ := s.Channel(v.ChannelID)
	affectedBeforChannel := ""

	if v.BeforeUpdate != nil {
		tmp, _ := s.Channel(v.BeforeUpdate.ChannelID)
		affectedBeforChannel = tmp.Name
	}

	if v.BeforeUpdate != nil {
		if v.BeforeUpdate.ChannelID == v.ChannelID {
			return
		}
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
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tartextchannel, err := s.GuildChannelCreate(v.GuildID, user.Username, discordgo.ChannelTypeGuildText)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tarkategory, err := s.GuildChannelCreate(v.GuildID, user.Username, discordgo.ChannelTypeGuildCategory)
		temprole, err := s.GuildRoleCreate(v.GuildID)

		s.GuildRoleEdit(v.GuildID, temprole.ID, "voicebotrole: "+temprole.ID, temprole.Color, temprole.Hoist, temprole.Permissions, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		m[targetchannel.ID] = GeneratedChannel{
			tarkategory.ID,
			targetchannel.ID,
			temprole.ID,
			tartextchannel.ID,
			0,
		}
		log.Printf("Channel erstellt!\n")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		manipulatePermissions(s, targetchannel.ID, tartextchannel.ID, tarkategory.ID, v.UserID, temprole.ID)

		dat0 := discordgo.ChannelEdit{
			Position: 4,
		}
		s.ChannelEditComplex(tarkategory.ID, &dat0)

		data := discordgo.ChannelEdit{
			ParentID: tarkategory.ID,
		}
		s.ChannelEditComplex(targetchannel.ID, &data)
		s.ChannelEditComplex(tartextchannel.ID, &data)

		err = s.GuildMemberMove(v.GuildID, v.UserID, &targetchannel.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	for key := range m {
		if v.ChannelID == key {
			tmpMap := m[key]
			tmpMap.current = tmpMap.current + 1
			m[key] = tmpMap
			s.GuildMemberRoleAdd(v.GuildID, v.UserID, tmpMap.temprole)
			log.Printf("%s - %d\n", key, m[key].current)
		}
	}

	if v.BeforeUpdate == nil {
		return
	}

	for key := range m {
		if v.BeforeUpdate.ChannelID == key {
			tmpMap := m[key]
			tmpMap.current = tmpMap.current - 1
			s.GuildMemberRoleRemove(v.GuildID, v.UserID, tmpMap.temprole)
			m[key] = tmpMap
			log.Printf("%s - %d\n", key, m[key].current)
			if m[key].current == 0 {
				s.GuildRoleDelete(v.GuildID, tmpMap.temprole)
				s.ChannelDelete(m[key].kategory)
				s.ChannelDelete(m[key].voicechannel)
				s.ChannelDelete(m[key].textchannel)
				delete(m, key)
				log.Printf("%s destroyed\n", key)
				return
			}

		}
	}

}
