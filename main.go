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

type SessionWrapper struct {
	*discord.Session
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INFO] No .env file found.")
	}

	if os.Getenv("CLIENT_TOKEN") == "" {
		log.Fatalln("[FATAL] CLIENT_TOKEN env. variable not specified.")
	}

	if os.Getenv("CLIENT_ID") == "" {
		log.Fatalln("[FATAL] CLIENT_ID env. variable not specified.")
	}
}

func main() {
	// Initialise new discord bot session
	s := SessionWrapper{initNewDiscordSession(os.Getenv("CLIENT_TOKEN"))}

	// Set intents
	s.Identify.Intents = discord.IntentsAll

	// Setup state tracking
	s.StateEnabled = true
	s.State.MaxMessageCount = 10

	s.addEventHandlers()

	s.addCommandHandlers()

	s.openConnection()

	s.registerCommands(os.Getenv("GUILD_ID"))

	// Close connection when process exits
	defer func() {
		s.closeConnection()
	}()

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	s.removeRegisteredCommands(os.Getenv("CLIENT_ID"), os.Getenv("GUILD_ID"))
}

func initNewDiscordSession(clientID string) *discord.Session {
	s, err := discord.New("Bot " + clientID)
	if err != nil {
		log.Fatalln("[FATAL] Could not create Discord Bot session.", err)
	}

	return s
}

func (s *SessionWrapper) addEventHandlers() {
	log.Println("[INFO] Adding event handlers...")
	for _, eh := range events.EventHandlers {
		s.AddHandler(eh)
	}
}

func (s *SessionWrapper) addCommandHandlers() {
	log.Println("[INFO] Adding command handlers...")
	s.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if c, ok := commands.SlashCommands[i.ApplicationCommandData().Name]; ok {
			c.CommandHandler(s, i)
		}
	})
}

func (s *SessionWrapper) openConnection() {
	err := s.Open()
	if err != nil {
		log.Fatalln("[FATAL] Could not open connection.", err)
	}
}

func (s *SessionWrapper) registerCommands(guildID string) {
	registerCommandsFlag := strings.ToLower(os.Getenv("REGISTER_COMMANDS"))
	if funk.Contains([]string{"1", "yes", "on", "true"}, registerCommandsFlag) {
		log.Println("[INFO] Registering commands...")
		var registeredCommands []*discord.ApplicationCommand
		for _, v := range commands.SlashCommands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v.CommandData)
			if err != nil {
				log.Panicf("[ERROR] Cannot create '%v' command: %v", v.CommandData.Name, err)
			}
			registeredCommands = append(registeredCommands, cmd)
		}
		log.Println("[INFO] Registered commands:")
		for _, v := range registeredCommands {
			log.Printf("[INFO] /%v", v.Name)
		}
	}
}

func (s *SessionWrapper) closeConnection() {
	log.Println("[INFO] Closing connection...")
	err := s.Close()
	if err != nil {
		log.Fatalln("[FATAL] Failed to close connection.", err)
	}
	log.Println("[INFO] Conneciton successfuly closed.")

}

func (s *SessionWrapper) removeRegisteredCommands(clientID string, guildID string) {
	removeCommandsFlag := strings.ToLower(os.Getenv("REMOVE_COMMANDS"))
	if funk.Contains([]string{"1", "yes", "on", "true"}, removeCommandsFlag) {
		log.Println("[INFO] Removing commands...")

		registeredCommands, err := s.ApplicationCommands(clientID, guildID)
		if err != nil {
			log.Panicf("[ERROR] Could not fetch registered commands: %v", err)
			return
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
			if err != nil {
				log.Panicf("[ERROR] Could not delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
