package eventhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"voicebot-discord/m/config"

	"github.com/bwmarrin/discordgo"
)

// GeneratedChannel is the Struct which contains all current informations of a
// Channel
type GeneratedChannel struct {
	Kategory     string
	Voicechannel string
	Textchannel  string
	Current      int
	Admin        string
}

var (
	m map[string]GeneratedChannel
)

func writeToJSON(foo map[string]GeneratedChannel) error {
	file, _ := os.Create("sqlite/Channel.json")
	defer file.Close()

	j, err := json.Marshal(foo)
	if err != nil {
		return fmt.Errorf("main: Error beim JSON Erstellen : %v", err)
	}
	fmt.Println(string(j))
	file.WriteString(string(j))
	return nil
}

// Init generates the map from an exisiting Json, or creates a new Map
func Init() {
	file, _ := os.ReadFile("sqlite/Channel.json")
	json.Unmarshal(file, &m)
	if m == nil {
		fmt.Println("Erstelle neue Liste...")
		m = map[string]GeneratedChannel{}
	}
}

func checkChannelID(v *discordgo.VoiceStateUpdate) error {
	if v.ChannelID == "foobar" {
		return nil
	}
	return nil
}

func giveUserPermission(s *discordgo.Session, textchannel string, user string) error {
	return s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1024, //allow
		0)    //deny
}
func takeUserPermission(s *discordgo.Session, textchannel string, user string) error {
	return s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		0,    //allow
		1024) //deny
}

// manipulatePermissions: Setzt die Permissions der neuen Channel richtig, sodass der Owner volle Rechte hat
// und nur Mitglieder im Channel/der Gruppe die Textchannel sehen
func manipulatePermissions(s *discordgo.Session, voicechannel string, textchannel string, kategorie string, user string) error {
	everybode := config.Everybody

	// Jeder in der ROlle darf den Textchannel sehen
	// Jeder ohne Rolle darf den Textchannel nicht sehen
	error := s.ChannelPermissionSet(
		textchannel,
		everybode,
		discordgo.PermissionOverwriteTypeRole,
		0,
		1024)
	if error != nil {
		return error
	}
	// Ersteller bekommt Admin rechte auf die Kategorie
	error = s.ChannelPermissionSet(
		kategorie,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)

	if error != nil {
		return error
	}
	// Ersteller auf Textchannel
	error = s.ChannelPermissionSet(
		textchannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)

	if error != nil {
		return error
	}
	// Ersteller auf Voicechannel
	error = s.ChannelPermissionSet(
		voicechannel,
		user,
		discordgo.PermissionOverwriteTypeMember,
		1040,
		0)
	return error
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
			return nil, errors.New("event ohne Channeländerung")
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

// createNewChannels: Handelt die komplette Interaktion
// beim Betreten des "main Channels"
func createNewChannels(s *discordgo.Session, v *discordgo.VoiceStateUpdate, user *discordgo.User) {

	tarkategoryname := user.Username
	tartextchannelname := user.Username
	targetchannelname := user.Username

	randomnames, err := GetRandomEntry()
	if err == nil {
		tarkategoryname = randomnames.Kategory
		tartextchannelname = randomnames.Textchannel
		targetchannelname = randomnames.Voicechannel
	}
	if err != nil {
		fmt.Printf("createNewChannels: %v", err)
	}

	targetchannel, err := s.GuildChannelCreate(
		v.GuildID,
		targetchannelname,
		discordgo.ChannelTypeGuildVoice)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tartextchannel, err := s.GuildChannelCreate(
		v.GuildID,
		tartextchannelname,
		discordgo.ChannelTypeGuildText)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tarkategory, err := s.GuildChannelCreate(
		v.GuildID,
		tarkategoryname,
		discordgo.ChannelTypeGuildCategory)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m[targetchannel.ID] = GeneratedChannel{
		tarkategory.ID,
		targetchannel.ID,
		tartextchannel.ID,
		0,
		v.UserID,
	}
	log.Printf("Channel erstellt!\n")

	manipulatePermissions(
		s,
		targetchannel.ID,
		tartextchannel.ID,
		tarkategory.ID,
		v.UserID)

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

// VoiceChannelCreate creates a new Voicechannel and moves the Creater there
// Calls all other Funcions to do the work and saves this in the JSON
func VoiceChannelCreate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	user, err := writeDownLog(s, v)
	defer writeToJSON(m)

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
			tmpMap.Current = tmpMap.Current + 1
			m[key] = tmpMap
			if user.ID != tmpMap.Admin {
				err := giveUserPermission(s, tmpMap.Textchannel, v.UserID)
				if err != nil {
					log.Println(err)
					return
				}
			}
			s.ChannelMessageSend(
				tmpMap.Textchannel,
				fmt.Sprintf("%s ist dem Channel beigetreten!", user.Username))
			log.Printf("%s - %d\n", key, m[key].Current)
		}
	}

	// Early Return bei initialen Join
	if v.BeforeUpdate == nil {
		return
	}

	// User verlässt Channel
	for key := range m {
		if v.BeforeUpdate.ChannelID == key {
			tmpMap := m[key]
			tmpMap.Current = tmpMap.Current - 1
			m[key] = tmpMap
			s.ChannelMessageSend(
				tmpMap.Textchannel,
				fmt.Sprintf("%s ist jetzt weg!", user.Username))
			log.Printf("%s - %d\n", key, m[key].Current)
			err := takeUserPermission(s, tmpMap.Textchannel, v.UserID)
			if err != nil {
				log.Println(err)
			}
			if m[key].Current == 0 {
				s.ChannelDelete(m[key].Voicechannel)
				s.ChannelDelete(m[key].Textchannel)
				s.ChannelDelete(m[key].Kategory)
				delete(m, key)
				log.Printf("%s destroyed\n", key)
				return
			}

		}
	}

}
