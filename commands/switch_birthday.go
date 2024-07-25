package commands

import (
	"log"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var SwitchBirthdayCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "switch-birthday",
		Description: "Switch birthday functionality on this server.",
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {
		guildData, err := utl.LoadGuildFromDBByID(i.GuildID)
		if err != nil {
			log.Printf("[ERROR] Could not load guild from database: %s", err)
		}

		if guildData.Birthday {
			guildData.Birthday = false
			utl.SaveGuildToDB(guildData)
			utl.RespondToInteractionCreateWithString(s, i, "Birthday functionality is now disabled.")
			return
		}

		guildData.Birthday = true
		utl.SaveGuildToDB(guildData)
		utl.RespondToInteractionCreateWithString(s, i, "Birthday functionality is now enabled.")
	},
}
