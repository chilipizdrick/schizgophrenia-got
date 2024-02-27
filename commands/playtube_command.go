package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	discord "github.com/bwmarrin/discordgo"
)

var PlayYoutubeCommand = SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "playtube",
		Description: "Plays Youtube video.",
		Options: []*discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionString,
				Name:        "query",
				Description: "YouTube search query",
				Required:    true,
			},
		},
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {

		const TEMP_AUDIO_FILE_LOCATION = "./temp/audio/youtube/"

		// log.Println("[TRACE] Running /playtube...")

		voiceChannelID, err := getInteractionVoiceChannelID(s, i)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}

		// Get YouTube query
		optionMap := getOptionMap(i)
		query, ok := optionMap["query"]
		if !ok {
			respondToInteractionCreateWithString(s, i, "No search query has been provided!")
			return
		}

		// Create direcrtory structure if it does not exist
		_, err = os.Stat(TEMP_AUDIO_FILE_LOCATION)
		if os.IsNotExist(err) {
			err = os.MkdirAll(TEMP_AUDIO_FILE_LOCATION, os.ModePerm)
			if err != nil {
				log.Printf("[ERROR] %v", err)
				respondToInteractionCreateWithString(s, i, "Could not download youtube audio!")
				return
			}
		}

		// Download youtube video
		log.Printf("[TRACE] Downloading youtube audio...")
		filepath := TEMP_AUDIO_FILE_LOCATION + "temp_audio.opus"
		err = downloadVideoWithYTDL(filepath, query.StringValue())
		if err != nil {
			log.Printf("[ERROR] %v", err)
			respondToInteractionCreateWithString(s, i, "Could not download youtube audio!")
			return
		}

		audioBuffer := make([][]byte, 0)

		log.Printf("[TRACE] Reading opus audio ...")
		err = loadOpusFile(filepath, &audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			respondToInteractionCreateWithString(s, i, "Could not load audio file!")
			return
		}

		err = playAudio(s, i.GuildID, voiceChannelID, audioBuffer)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			respondToInteractionCreateWithString(s, i, "Could not connect to voive channel!")
			return
		}
	},
}

func downloadVideoWithYTDL(filepath, query string) error {

	command := "youtube-dl"
	command += fmt.Sprintf(" --output \"%v\"", filepath)
	command += " --extract-audio"
	command += " --audio-format opus"
	// command += " --geo-bypass"

	command += fmt.Sprintf(" \"ytsearch1:%v\"", query)

	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run() // waits until the commands runs and finishes
	return err
}

// func parseYouTubeLink(link string) (string, error) {
// 	watchIndex := strings.Index(link, "watch?v=")
// 	if watchIndex != -1 {
// 		videoIDIndex := watchIndex + len("watch?v=")
// 		return link[videoIDIndex : videoIDIndex+11], nil
// 	}

// 	youtubeIndex := strings.Index(link, "youtu.be/")
// 	if youtubeIndex != -1 {
// 		videoIDIndex := youtubeIndex + len("youtu.be/")
// 		return link[videoIDIndex : videoIDIndex+11], nil
// 	}

// 	return "", errors.New("parse error: invalid youtube link")
// }
