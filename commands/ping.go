package commands

import "github.com/bwmarrin/discordgo"

var PingCommand = SlashCommand{
	CommandData: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with \"Pong!\"",
	},
	CommandHandler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
	},
}
