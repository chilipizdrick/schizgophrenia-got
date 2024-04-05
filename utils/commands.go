package utils

import (
	"errors"
	"log"

	discord "github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	CommandData    *discord.ApplicationCommand
	CommandHandler func(*discord.Session, *discord.InteractionCreate)
}

func GetOptionMap(i *discord.InteractionCreate) map[string]*discord.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discord.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func RespondToInteractionCreateWithString(s *discord.Session, i *discord.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: message,
		},
	})
}

func EditResponseWithString(s *discord.Session, i *discord.InteractionCreate, message string) {
	s.InteractionResponseEdit(i.Interaction, &discord.WebhookEdit{
		Content: &message,
	})
}

func GetInteractionVoiceChannelID(s *discord.Session, i *discord.InteractionCreate) (string, error) {
	voiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return "", errors.New("could not get user's voice state")
	}

	voiceChannelID := voiceState.ChannelID
	if voiceChannelID == "" {
		return "", errors.New("user not in voice channel")
	}

	return voiceChannelID, nil
}

func GenericVoiceCommandHandler(filepath string) func(s *discord.Session, i *discord.InteractionCreate) {
	return func(s *discord.Session, i *discord.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseDeferredChannelMessageWithSource,
		})
		defer s.InteractionResponseDelete(i.Interaction)

		// Return if user is bot
		if i.Member.User.Bot {
			return
		}

		voiceChannelID, err := GetInteractionVoiceChannelID(s, i)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			EditResponseWithString(s, i, "User should be in voice channel.")
			return
		}

		audioBuffer := make([][]byte, 0)

		err = LoadOpusFile(filepath, &audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			EditResponseWithString(s, i, "Could not load audio file!")
			return
		}

		err = PlayAudio(s, i.GuildID, voiceChannelID, audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			// editResponseWithString(s, i, "Could not connect to voice channel or disconnect from it!")
			return
		}
	}
}

func GenericRandomVoiceCommandHandler(dirpath string) func(s *discord.Session, i *discord.InteractionCreate) {
	return func(s *discord.Session, i *discord.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseDeferredChannelMessageWithSource,
		})
		defer s.InteractionResponseDelete(i.Interaction)

		// Return if user is bot
		if i.Member.User.Bot {
			return
		}

		voiceChannelID, err := GetInteractionVoiceChannelID(s, i)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			EditResponseWithString(s, i, "User should be in voice channel.")
			return
		}

		audioBuffer := make([][]byte, 0)

		filepath, err := PickRandomFileFromDirectory(dirpath)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			EditResponseWithString(s, i, "Could not randomly pick audio file path!")
			return
		}
		err = LoadOpusFile(filepath, &audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			EditResponseWithString(s, i, "Could not load audio file!")
			return
		}

		err = PlayAudio(s, i.GuildID, voiceChannelID, audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			// editResponseWithString(s, i, "Could not connect to voice channel or disconnect from it!")
			return
		}
	}
}

func GenericVoiceCommand(name, description, filepath string) SlashCommand {
	return SlashCommand{
		CommandData: &discord.ApplicationCommand{
			Name:        name,
			Description: description,
		},
		CommandHandler: GenericVoiceCommandHandler(filepath),
	}
}

func GenericRandomVoiceCommand(name, descriprion, dirpath string) SlashCommand {
	return SlashCommand{
		CommandData: &discord.ApplicationCommand{
			Name:        name,
			Description: descriprion,
		},
		CommandHandler: GenericRandomVoiceCommandHandler(dirpath),
	}
}
