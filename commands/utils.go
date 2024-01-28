package commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getOptionMap(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func respondToInteractionWithString(s *discordgo.Session, i *discordgo.Interaction, message string) {

	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func loadAudioFile(filepath string, buffer *[][]byte) error {

	var opuslen int16

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	for {
		err := binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from file :", err)
			return err
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from file :", err)
			return err
		}

		*buffer = append(*buffer, inBuf)
	}
}

func playSound(s *discordgo.Session, guildID string, channelID string, buffer [][]byte) (err error) {

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	vc.Disconnect()

	return nil
}
