package events

import (
	"fmt"
	"log"
	"os"
	"regexp"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

func MessageCreateHandler(s *discord.Session, e *discord.MessageCreate) {
	guildData, err := utl.LoadGuildFromDBByID(e.GuildID)
	if err != nil {
		log.Printf("[ERROR] Could not load guild from databse by id: %e", err)
	}

	if guildData.DMOnMention {
		if e.Author == nil {
			return
		}

		if e.Author.ID == "" || e.Author.ID == os.Getenv("CLIENT_ID") {
			return
		}

		match, err := regexp.MatchString(".*<@[0-9]{18}>.*", e.Message.Content)
		if err != nil {
			log.Printf("[ERROR] Could not match string: %e", err)
			return
		}
		if !match {
			return
		}

		dmChannel, err := s.UserChannelCreate(e.Author.ID)
		if err != nil {
			log.Printf("[ERROR] Could not open dm channel: %e", err)
			return
		}
		guild, err := s.State.Guild(e.GuildID)
		if err != nil {
			log.Printf("[ERROR] Could not fetch guild by id: %e", err)
			return
		}
		mentionChannel, err := s.State.Channel(e.ChannelID)
		if err != nil {
			log.Printf("[ERROR] Could not fetch channel by id: %e", err)
			return
		}
		messageContents := fmt.Sprintf("You (%s) have been mentioned if channel %v of server %s.", e.Author.Mention(), mentionChannel.Name, guild.Name)
		_, err = s.ChannelMessageSend(dmChannel.Mention(), messageContents)
		if err != nil {
			log.Printf("[ERROR] Could not send a dm: %e", err)
			return
		}
	}

}
