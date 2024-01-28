package events

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyEvent(s *discordgo.Session, e *discordgo.Ready) {
	// log.Printf("[TRACE] %v", s.State)
	logMessage := fmt.Sprintf("[INFO] Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	log.Println(logMessage)
}
