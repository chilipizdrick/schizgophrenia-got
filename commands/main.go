package commands

import (
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var (
	SlashCommands = map[string]utl.SlashCommand{
		PingCommand.CommandData.Name:        PingCommand,
		ColorCommand.CommandData.Name:       ColorCommand,
		RemoveColorCommand.CommandData.Name: RemoveColorCommand,

		// Generic voice commands
		"pipe":   utl.GenericVoiceCommand("pipe", "Plays metal pipe sound", "./assets/audio/voice/pipe.ogg"),
		"fontan": utl.GenericVoiceCommand("fontan", "Plays \"Chocoladniy Fontan\"", "./assets/audio/voice/fontan.ogg"),
		"women":  utl.GenericVoiceCommand("women", "Plays women", "./assets/audio/voice/women.ogg"),
	}
)
