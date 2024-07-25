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
			log.Printf("[ERROR] Error fetching guild roles by guildID. %s", err)
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
				log.Printf("[ERROR] Error creating new personal user's role. %s", err)
				utl.RespondToInteractionCreateWithString(s, i, "Could not create user's personal role.")
				return
			}

			err = s.GuildMemberRoleAdd(i.GuildID, userID, personalRole.ID)
			if err != nil {
				utl.RespondToInteractionCreateWithString(s, i, "Could not give newly created personal role to user.")
				log.Printf("[ERROR] Could not give newly created personal role to user. %s", err)
				return
			}
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

		// Check if role exists, but user does not have it
		// (may happen when user leaves and rejoins the server)
		member, err := s.State.Member(i.GuildID, userID)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "Could not fetch user by ID.")
			log.Printf("[ERROR] Error fetching member by ID. %s", err)
			return
		}
		if !funk.Contains(member.Roles, personalRole.ID) {
			err = s.GuildMemberRoleAdd(i.GuildID, userID, personalRole.ID)
			if err != nil {
				utl.RespondToInteractionCreateWithString(s, i, "Could not give personal role to user.")
				log.Printf("[ERROR] Could not give personal role to user. %s", err)
				return
			}
		}

		optionMap := utl.GetOptionMap(i)
		colorField := optionMap["color"].StringValue()
		cleanedHexColor := strings.Replace(colorField, "#", "", -1)
		uIntColor, err := strconv.ParseUint(cleanedHexColor, 16, 64)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "Invalid color in HEX format has been provided.")
			log.Printf("[ERROR] Error parsing hexidecimal color. %s", err)
			return
		}

		intColor := int(uIntColor)

		s.GuildRoleEdit(i.GuildID, personalRole.ID, &discord.RoleParams{
			Color: &intColor,
		})

		utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated %s's personal role color.", i.Member.Mention()))
	},
}
