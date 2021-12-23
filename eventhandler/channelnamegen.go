package eventhandler

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"voicebot-discord/m/database"

	"github.com/bwmarrin/discordgo"
)

func GenerateNewEntry(s *discordgo.Session, v *discordgo.MessageCreate) {
	content := strings.SplitAfter(v.Content, " ")
	// Ignore everythin not starting with !add
	//s.ChannelMessageSend(v.ChannelID, content[0])

	if !(strings.HasPrefix(content[0], "!add")) {
		return
	}

	if len(content) != 4 {
		s.ChannelMessageSend(v.ChannelID, "Bitte mit \"!add $KATEGORIENAME $TEXTCHANNELNAME $VOICECHANNELNAME\" aufrufen. Leerzeichen kann man statt dessen mit _ Schreiben")
		return
	}

	for i := range content {
		content[i] = strings.ReplaceAll(content[i], "_", " ")
	}

	newentry := database.Entry{
		content[1],
		content[2],
		content[3],
		v.Author.Username,
	}
	database.AddEntry(newentry.Kategory, newentry.Textchannel, newentry.Voicechannel, newentry.Creator)

	//entries[id] = newentry

	output = "Kategorie: " + newentry.Kategorie + "\nTextchannel: " + newentry.Textchannel + "\nVoicechannel: " + newentry.Voicechannel + "\n By " + newentry.Creator

	s.ChannelMessageSend(v.ChannelID, output)
}

func GetRandomEntry() (database.Entry, error) {

	min := 0
	entries := database.GetAllEntrys()
	max := len(entries)

	fmt.Println(entries)

	if max == 0 {
		return database.Entry{"", "", "", ""}, errors.New("Keine Daten vorhanden! Füge vorher neue Entries hinzu!")
	}

	randomentry := rand.Intn(max-min) + min

	fmt.Print(entries[randomentry])
	return entries[randomentry], nil
}
