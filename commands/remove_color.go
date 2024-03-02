package commands

import (
	"fmt"
	"log"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var RemoveColorCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "remove-color",
		Description: "Removes user's personal color role",
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {

		userID := i.Member.User.ID

		guildRoles, err := s.GuildRoles(i.GuildID)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "Could not get server roles.")
			log.Printf("[ERROR] Error fetching guild roles by guildID. %v", err)
			return
		}

		var personalRole *discord.Role
		for _, role := range guildRoles {
			if role.Name == userID {
				personalRole = role
				break
			}
		}
		if personalRole == nil {
			utl.RespondToInteractionCreateWithString(s, i, "User has no personal role.")
			return
		}

		s.GuildRoleDelete(i.GuildID, personalRole.ID)

		utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully deleted %v's personal color role.", i.Member.Mention()))
	},
}
