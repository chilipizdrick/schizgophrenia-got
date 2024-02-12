package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var RemoveColorCommand = SlashCommand{
	CommandData: &discordgo.ApplicationCommand{
		Name:        "remove-color",
		Description: "Removes user's personal color role",
	},
	CommandHandler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		userID := i.Member.User.ID

		guildRoles, err := s.GuildRoles(i.GuildID)
		if err != nil {
			respondToInteractionCreateWithString(s, i, "Could not get server roles.")
			log.Printf("[ERROR] Error fetching guild roles by guildID. %v", err)
			return
		}

		var personalRole *discordgo.Role
		for _, role := range guildRoles {
			if role.Name == userID {
				personalRole = role
				break
			}
		}
		if personalRole == nil {
			respondToInteractionCreateWithString(s, i, "User has no personal role.")
			return
		}

		s.GuildRoleDelete(i.GuildID, personalRole.ID)

		respondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully deleted %v's personal color role.", i.Member.Mention()))
	},
}
