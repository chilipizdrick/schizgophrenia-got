package commands

import "github.com/bwmarrin/discordgo"

type SlashCommand struct {
	CommandData    *discordgo.ApplicationCommand
	CommandHandler func(*discordgo.Session, *discordgo.InteractionCreate)
}

var (
	SlashCommands = map[string]SlashCommand{
		PingCommand.CommandData.Name:        PingCommand,
		PlayYoutubeCommand.CommandData.Name: PlayYoutubeCommand,
		// Generic voice commands
		"pipe":   genericVoiceCommand("pipe", "Plays metal pipe sound", "./assets/audio/voice/pipe.ogg"),
		"fontan": genericVoiceCommand("fontan", "Plays \"Chocoladniy Fontan\"", "./assets/audio/voice/fontan.ogg"),
	}
)
