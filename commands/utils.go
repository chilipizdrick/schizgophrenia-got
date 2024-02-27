package commands

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	discord "github.com/bwmarrin/discordgo"
	"github.com/pion/opus/pkg/oggreader"
)

func getOptionMap(i *discord.InteractionCreate) map[string]*discord.ApplicationCommandInteractionDataOption {

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discord.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func respondToInteractionCreateWithString(s *discord.Session, i *discord.InteractionCreate, message string) {

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: message,
		},
	})
}

func editResponseWithString(s *discord.Session, i *discord.InteractionCreate, message string) {

	s.InteractionResponseEdit(i.Interaction, &discord.WebhookEdit{
		Content: &message,
	})
}

func getInteractionVoiceChannelID(s *discord.Session, i *discord.InteractionCreate) (string, error) {

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

func loadOpusFile(filepath string, buffer *[][]byte) error {

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	defer file.Close()

	reader, _, err := oggreader.NewWith(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	for {
		inBuf, _, err := reader.ParseNextPage()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		*buffer = append(*buffer, inBuf...)
	}
}

func playAudio(s *discord.Session, guildID string, channelID string, buffer [][]byte) (err error) {

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)

	for _, frame := range buffer {
		vc.OpusSend <- frame
	}

	vc.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	err = vc.Disconnect()
	if err != nil {
		return fmt.Errorf("error on voiceChannnel.Disconect: %w", err)
	}

	return nil
}

func genericVoiceCommandHandler(filepath string) func(s *discord.Session, i *discord.InteractionCreate) {

	return func(s *discord.Session, i *discord.InteractionCreate) {

		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseDeferredChannelMessageWithSource,
		})

		defer s.InteractionResponseDelete(i.Interaction)

		voiceChannelID, err := getInteractionVoiceChannelID(s, i)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		audioBuffer := make([][]byte, 0)

		err = loadOpusFile(filepath, &audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			editResponseWithString(s, i, "Could not load audio file!")
			return
		}

		err = playAudio(s, i.GuildID, voiceChannelID, audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			// editResponseWithString(s, i, "Could not connect to voice channel or disconnect from it!")
			return
		}
	}
}

func genericVoiceCommand(name, description, filepath string) SlashCommand {

	return SlashCommand{
		CommandData: &discord.ApplicationCommand{
			Name:        name,
			Description: description,
		},
		CommandHandler: genericVoiceCommandHandler(filepath),
	}
}
