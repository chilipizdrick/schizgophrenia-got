package main

import (
	"log"
	"os"
	"testing"

	discord "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func TestBasicFunctionality(t *testing.T) {

	// Load and check for env. variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INFO] No .env file found.")
	}
	if os.Getenv("TESTING_CLIENT_TOKEN") == "" {
		log.Fatalln("[FATAL] TESTING_CLIENT_TOKEN env. variable not specified.")
	}
	if os.Getenv("TESTING_CLIENT_ID") == "" {
		log.Fatalln("[FATAL] TESTING_CLIENT_ID env. variable not specified.")
	}

	s := SessionWrapper{initNewDiscordSession(os.Getenv("TESTING_CLIENT_TOKEN"))}

	s.Identify.Intents = discord.IntentsAll

	s.StateEnabled = true
	s.State.MaxMessageCount = 10

	s.addEventHandlers()

	s.addCommandHandlers()

	s.openConnection()

	s.registerCommands(os.Getenv("TESTING_GUILD_ID"))

	s.removeRegisteredCommands(os.Getenv("TESTING_CLIENT_ID"), os.Getenv("TESTING_GUILD_ID"))

	s.closeConnection()
}
