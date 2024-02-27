package commands

import discord "github.com/bwmarrin/discordgo"

var PingCommand = SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with \"Pong!\"",
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Content: "Pong!",
			},
		})
	},
}
