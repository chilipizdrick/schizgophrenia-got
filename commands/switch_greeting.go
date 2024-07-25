package commands

import (
	"log"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var SwitchGreetingCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "switch-greeting",
		Description: "Switch greeting functionality on this server.",
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {
		guildData, err := utl.LoadGuildFromDBByID(i.GuildID)
		if err != nil {
			log.Printf("[ERROR] Could not load guild from database: %s", err)
		}

		if guildData.Greeting {
			guildData.Greeting = false
			utl.SaveGuildToDB(guildData)
			utl.RespondToInteractionCreateWithString(s, i, "Greeting functionality is now disabled.")
			return
		}

		guildData.Greeting = true
		utl.SaveGuildToDB(guildData)
		utl.RespondToInteractionCreateWithString(s, i, "Greeting functionality is now enabled.")
	},
}
