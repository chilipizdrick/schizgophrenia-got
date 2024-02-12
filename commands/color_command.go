package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thoas/go-funk"
)

var ColorCommand = SlashCommand{
	CommandData: &discordgo.ApplicationCommand{
		Name:        "color",
		Description: "Changes user's personal role color",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "color",
				Description: "New role color in HEX format",
				Required:    true,
			},
		},
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
		if !funk.Contains(funk.Map(guildRoles, func(role *discordgo.Role) string {
			return role.Name
		}), userID) {
			permissions := int64(0)
			hoist := false
			mentionable := false

			personalRole, err = s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
				Name:        userID,
				Permissions: &permissions,
				Hoist:       &hoist,
				Mentionable: &mentionable,
			})
			if err != nil {
				log.Printf("[ERROR] Error creating new personal user's role. %v", err)
				respondToInteractionCreateWithString(s, i, "Could not create user's personal role.")
				return
			}

			s.GuildMemberRoleAdd(i.GuildID, userID, personalRole.ID)
		}

		// If role has not been created then fetch it
		if personalRole == nil {
			for _, role := range guildRoles {
				if role.Name == userID {
					personalRole = role
					break
				}
			}
		}

		optionMap := getOptionMap(i)
		colorField := optionMap["color"].StringValue()
		cleanedHexColor := strings.Replace(colorField, "#", "", -1)
		uIntColor, err := strconv.ParseUint(cleanedHexColor, 16, 64)
		if err != nil {
			respondToInteractionCreateWithString(s, i, "Invalid color in HEX format has been provided.")
			log.Printf("[ERROR] Error parsing hexidecimal color. %v", err)
			return
		}

		intColor := int(uIntColor)

		s.GuildRoleEdit(i.GuildID, personalRole.ID, &discordgo.RoleParams{
			Color: &intColor,
		})

		respondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated %v's personal role color.", i.Member.Mention()))
	},
}
