package commands

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/lithdew/nicehttp"
	"github.com/lithdew/youtube"
)

var PlayYoutubeCommand = SlashCommand{
	CommandData: &discordgo.ApplicationCommand{
		Name:        "playtube",
		Description: "Plays Youtube video.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "YouTube search query",
				Required:    true,
			},
		},
	},
	CommandHandler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		const TEMP_AUDIO_FILE_LOCATION = "./temp/audio/youtube/"

		log.Println("[TRACE] Running /playtube...")

		// Fetch member's voice channel (requires s.StateEnabled to be true)
		// log.Printf("[TRACE] Guild ID: %v", i.GuildID)
		// log.Printf("[TRACE] User ID: %v", i.Member.User.ID)
		voiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not get user's current voice channel data!")
			return
		}
		voiceChannelID := voiceState.ChannelID
		// log.Printf("[TRACE] Voice channel ID: %v", voiceChannelID)
		if voiceChannelID == "" {
			respondToInteractionWithString(s, i.Interaction, "This command can only be used in a voice channel!")
			return
		}

		// Get YouTube query
		optionMap := getOptionMap(i)
		query, ok := optionMap["query"]
		if !ok {
			respondToInteractionWithString(s, i.Interaction, "No search query has been provided!")
			return
		}
		// log.Printf("[TRACE] YouTube query: %v", link.StringValue())

		// Search youtube with provided query
		results, err := youtube.Search(query.StringValue(), 0)
		if err != nil {
			log.Printf("[ERROR] Could not load youtube search results. %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not load youtube video!")
			return
		}

		if len(results.Items) == 0 {
			respondToInteractionWithString(s, i.Interaction, "Found no videos with provided query!")
			return
		}

		details := results.Items[0]
		log.Printf("[TRACE] Video details: %v", details)
		player, err := youtube.Load(details.ID)
		if err != nil {
			log.Printf("[ERROR] Could not load video player. %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not download youtube audio!")
			return
		}

		stream, ok := player.SourceFormats().AudioOnly().BestAudio()
		if !ok {
			log.Printf("[ERROR] Could not get audio stream. %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not download youtube audio!")
			return
		}

		// filename := "temp_youtube_audio." + stream.FileExtension()

		url, err := player.ResolveURL(stream)
		if err != nil {
			log.Printf("[ERROR] Could not resolve audio stream url. %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not download youtube audio!")
			return
		}

		// Create direcrtory structure if it does not exist
		_, err = os.Stat(TEMP_AUDIO_FILE_LOCATION)
		if os.IsNotExist(err) {
			err = os.MkdirAll(TEMP_AUDIO_FILE_LOCATION, os.ModePerm)
			if err != nil {
				respondToInteractionWithString(s, i.Interaction, "Could not download youtube audio!")
				return
			}
		}

		filename := "temp_youtube_audio." + stream.FileExtension()
		filepath := TEMP_AUDIO_FILE_LOCATION + filename

		err = nicehttp.DownloadFile(filepath, url)
		if err != nil {
			log.Printf("[ERROR] Could not download youtube video. %v", err)
			respondToInteractionWithString(s, i.Interaction, "Could not download youtube audio!")
			return
		}

		audioBuffer := make([][]byte, 0)

		err = loadAudioFile(filepath, &audioBuffer)
		if err != nil {
			respondToInteractionWithString(s, i.Interaction, "Could not load audio file!")
			return
		}

		playSound(s, i.GuildID, voiceChannelID, audioBuffer)
	},
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

// 	return "", errors.New("parse error: Invalid youtube link")
// }
