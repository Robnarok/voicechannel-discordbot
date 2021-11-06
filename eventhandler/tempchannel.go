package eventhandler

import (
	"errors"
	"fmt"
	"log"
	"voicebot-discord/m/config"

	"github.com/bwmarrin/discordgo"
)

type GeneratedChannel struct {
	kategory     string
	voicechannel string
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

func giveUserPermission(s *discordgo.Session, textchannel string, user string) {
	s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1024,
		0)
}
func takeUserPermission(s *discordgo.Session, textchannel string, user string) {
	s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		0,
		1024)
}

// manipulatePermissions: Setzt die Permissions der neuen Channel richtig, sodass der Owner volle Rechte hat
// und nur Mitglieder im Channel/der Gruppe die Textchannel sehen
func manipulatePermissions(s *discordgo.Session, voicechannel string, textchannel string, kategorie string, user string) {
	everybode := "835121335851155466"

	// Jeder in der ROlle darf den Textchannel sehen
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

// writeDownLog: Schreibt die Logs, welches Event passiert ist... Sorgt auch dafür, dass es
// später nicht zu NUllpointern kommt... Das ist nicht ganz so sauber, muss dringed refactored werden
func writeDownLog(s *discordgo.Session, v *discordgo.VoiceStateUpdate) (*discordgo.User, error) {
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
			return nil, errors.New("Crash")
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
	return user, nil
}

// createNewChannels: Handelt die komplette Interaktion beim Betreten des "main Channels"
func createNewChannels(s *discordgo.Session, v *discordgo.VoiceStateUpdate, user *discordgo.User) {
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
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m[targetchannel.ID] = GeneratedChannel{
		tarkategory.ID,
		targetchannel.ID,
		tartextchannel.ID,
		0,
	}
	log.Printf("Channel erstellt!\n")

	manipulatePermissions(s, targetchannel.ID, tartextchannel.ID, tarkategory.ID, v.UserID)

	dat0 := discordgo.ChannelEdit{
		Position: 3,
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

func VoiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	user, err := writeDownLog(s, v)

	if err != nil {
		log.Println(err.Error())
		return
	}
	if v.ChannelID == config.Masterchannel {
		createNewChannels(s, v, user)
	}

	// User tritt einem Voicechannel bei
	for key := range m {
		if v.ChannelID == key {
			tmpMap := m[key]
			tmpMap.current = tmpMap.current + 1
			m[key] = tmpMap
			giveUserPermission(s, tmpMap.textchannel, v.UserID)
			s.ChannelMessageSend(tmpMap.textchannel, fmt.Sprintf("%s ist dem Channel beigetretten!", user.Username))
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
			m[key] = tmpMap
			log.Printf("%s - %d\n", key, m[key].current)
			if m[key].current == 0 {
				takeUserPermission(s, tmpMap.textchannel, v.UserID)
				s.ChannelMessageSend(tmpMap.textchannel, fmt.Sprintf("%s ist jetzt weg!", user.Username))
				s.ChannelDelete(m[key].voicechannel)
				s.ChannelDelete(m[key].textchannel)
				s.ChannelDelete(m[key].kategory)
				delete(m, key)
				log.Printf("%s destroyed\n", key)
				return
			}

		}
	}

}
