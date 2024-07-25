package commands

import (
	"log"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var SwitchDMOnMentionCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "switch-dm-on-mention",
		Description: "[NOT CURRNETLY WORKING] Switch DM on mention functionality on this server.",
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {
		guildData, err := utl.LoadGuildFromDBByID(i.GuildID)
		if err != nil {
			log.Printf("[ERROR] Could not load guild from database: %s", err)
		}

		if guildData.DMOnMention {
			guildData.DMOnMention = false
			utl.SaveGuildToDB(guildData)
			utl.RespondToInteractionCreateWithString(s, i, "DM on mention functionality is now disabled.")
			return
		}

		guildData.DMOnMention = true
		utl.SaveGuildToDB(guildData)
		utl.RespondToInteractionCreateWithString(s, i, "DM on mention functionality is now enabled.")
	},
}
