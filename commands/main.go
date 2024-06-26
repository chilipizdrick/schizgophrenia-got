package commands

import (
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

var (
	SlashCommands = map[string]utl.SlashCommand{
		PingCommand.CommandData.Name:              PingCommand,
		ColorCommand.CommandData.Name:             ColorCommand,
		RemoveColorCommand.CommandData.Name:       RemoveColorCommand,
		RegisterBirthdayCommand.CommandData.Name:  RegisterBirthdayCommand,
		SwitchGreetingCommand.CommandData.Name:    SwitchGreetingCommand,
		SwitchBirthdayCommand.CommandData.Name:    SwitchBirthdayCommand,
		SwitchDMOnMentionCommand.CommandData.Name: SwitchDMOnMentionCommand,

		// Generic voice commands
		"pipe":     utl.GenericVoiceCommand("pipe", "Plays metal pipe sound", "./assets/audio/voice/pipe.ogg"),
		"fontan":   utl.GenericVoiceCommand("fontan", "Plays \"Chocoladniy Fontan\"", "./assets/audio/voice/fontan.ogg"),
		"women":    utl.GenericVoiceCommand("women", "Plays women", "./assets/audio/voice/women.ogg"),
		"oblivion": utl.GenericVoiceCommand("oblivion", "Plays Oblivion NPC theme", "./assets/audio/voice/oblivion.ogg"),
		"cave":     utl.GenericRandomVoiceCommand("cave", "Plays random minecraft cave sound", "./assets/audio/voice/cave/"),
	}
)
