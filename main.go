package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"github.com/chilipizdrick/schizgophrenia-got/commands"
	"github.com/chilipizdrick/schizgophrenia-got/events"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("[FATAL] Could not load .env file.")
	}

	if os.Getenv("CLIENT_TOKEN") == "" {
		log.Fatalln("[FATAL] Could not read CLIENT_TOKEN from the .env file.")
	}
}

func main() {

	// Initialise new discord bot session
	s, err := discord.New("Bot " + os.Getenv("CLIENT_TOKEN"))
	if err != nil {
		log.Fatalln("[FATAL] Could not create Discord Bot session.", err)
	}

	// Set intents
	s.Identify.Intents = discord.IntentsAll

	// Setup state tracking
	s.StateEnabled = true
	s.State.MaxMessageCount = 10

	// Add event handlers
	log.Println("[INFO] Adding event handlers...")
	for _, eh := range events.EventHandlers {
		s.AddHandler(eh)
	}

	// Add command handlers
	log.Println("[INFO] Adding command handlers...")
	s.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if c, ok := commands.SlashCommands[i.ApplicationCommandData().Name]; ok {
			c.CommandHandler(s, i)
		}
	})

	// Open connection
	err = s.Open()
	if err != nil {
		log.Fatalln("[FATAL] Could not open connection.", err)
	}

	// Register commands
	log.Println("[INFO] Registering commands...")
	var registeredCommands []*discord.ApplicationCommand
	guildId := os.Getenv("GUILD_ID")
	for _, v := range commands.SlashCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v.CommandData)
		if err != nil {
			log.Panicf("[ERROR] Cannot create '%v' command: %v", v.CommandData.Name, err)
		}
		registeredCommands = append(registeredCommands, cmd)
	}
	log.Println("[INFO] Registered commands:")
	for _, v := range registeredCommands {
		log.Printf("[INFO] /%v", v.Name)
	}

	// Close connection when process exits
	defer func() {
		log.Println("[INFO] Closing connection...")
		err := s.Close()
		if err != nil {
			log.Fatalln("[FATAL] Failed to close connection.", err)
		}
		log.Println("[INFO] Conneciton successfuly closed.")
	}()

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Remove registered commands
	removeCommandsFlag := strings.ToLower(os.Getenv("REMOVE_COMMANDS"))
	if funk.Contains([]string{"1", "yes", "on", "true"}, removeCommandsFlag) {
		log.Println("[INFO] Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildId, v.ID)
			if err != nil {
				log.Panicf("[ERROR] Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
