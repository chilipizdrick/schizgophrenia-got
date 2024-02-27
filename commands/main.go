package commands

import discord "github.com/bwmarrin/discordgo"

type SlashCommand struct {
	CommandData    *discord.ApplicationCommand
	CommandHandler func(*discord.Session, *discord.InteractionCreate)
}

var (
	SlashCommands = map[string]SlashCommand{
		PingCommand.CommandData.Name:        PingCommand,
		ColorCommand.CommandData.Name:       ColorCommand,
		RemoveColorCommand.CommandData.Name: RemoveColorCommand,
		PlayYoutubeCommand.CommandData.Name: PlayYoutubeCommand,
		// Generic voice commands
		"pipe":   genericVoiceCommand("pipe", "Plays metal pipe sound", "./assets/audio/voice/pipe.ogg"),
		"fontan": genericVoiceCommand("fontan", "Plays \"Chocoladniy Fontan\"", "./assets/audio/voice/fontan.ogg"),
		"women":  genericVoiceCommand("women", "Plays women", "./assets/audio/voice/women.ogg"),
	}
)
