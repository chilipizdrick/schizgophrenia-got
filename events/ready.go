package events

import (
	"fmt"
	"log"

	discord "github.com/bwmarrin/discordgo"
)

func ReadyEvent(s *discord.Session, e *discord.Ready) {
	// log.Printf("[TRACE] %v", s.State)
	logMessage := fmt.Sprintf("[INFO] Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	log.Println(logMessage)
}
