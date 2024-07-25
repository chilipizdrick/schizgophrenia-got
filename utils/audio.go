package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	discord "github.com/bwmarrin/discordgo"
	"github.com/pion/opus/pkg/oggreader"
)

func LoadOpusFile(filepath string, buffer *[][]byte) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}

	defer file.Close()

	reader, _, err := oggreader.NewWith(file)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}

	for {
		inBuf, _, err := reader.ParseNextPage()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("error reading file: %s", err)
		}

		*buffer = append(*buffer, inBuf...)
	}
}

func PlayAudio(s *discord.Session, guildID string, channelID string, buffer [][]byte) error {
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
		return fmt.Errorf("error on voiceChannnel.Disconect: %s", err)
	}

	return nil
}
