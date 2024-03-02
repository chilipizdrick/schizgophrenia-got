package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
	"github.com/thoas/go-funk"
)

var ColorCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "color",
		Description: "Changes user's personal role color",
		Options: []*discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionString,
				Name:        "color",
				Description: "New role color in HEX format",
				Required:    true,
			},
		},
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
		if !funk.Contains(funk.Map(guildRoles, func(role *discord.Role) string {
			return role.Name
		}), userID) {
			permissions := int64(0)
			hoist := false
			mentionable := false

			personalRole, err = s.GuildRoleCreate(i.GuildID, &discord.RoleParams{
				Name:        userID,
				Permissions: &permissions,
				Hoist:       &hoist,
				Mentionable: &mentionable,
			})
			if err != nil {
				log.Printf("[ERROR] Error creating new personal user's role. %v", err)
				utl.RespondToInteractionCreateWithString(s, i, "Could not create user's personal role.")
				return
			}

			s.GuildMemberRoleAdd(i.GuildID, userID, personalRole.ID)
		}

		// If role has not been created (already existed) then fetch it
		if personalRole == nil {
			for _, role := range guildRoles {
				if role.Name == userID {
					personalRole = role
					break
				}
			}
		}

		optionMap := utl.GetOptionMap(i)
		colorField := optionMap["color"].StringValue()
		cleanedHexColor := strings.Replace(colorField, "#", "", -1)
		uIntColor, err := strconv.ParseUint(cleanedHexColor, 16, 64)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "Invalid color in HEX format has been provided.")
			log.Printf("[ERROR] Error parsing hexidecimal color. %v", err)
			return
		}

		intColor := int(uIntColor)

		s.GuildRoleEdit(i.GuildID, personalRole.ID, &discord.RoleParams{
			Color: &intColor,
		})

		utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated %v's personal role color.", i.Member.Mention()))
	},
}
