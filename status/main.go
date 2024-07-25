package status

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	discord "github.com/bwmarrin/discordgo"
	"github.com/dreamscached/minequery/v2"
)

const QUERY_TIMEOUT = 10 * time.Second

func MinecraftServerCheckRoutine(s *discord.Session) {
	minecraftServerIp := os.Getenv("MINECRAFT_SERVER_IP")
	minecraftServerPortString := os.Getenv("MINECRAFT_SERVER_PORT")

	if minecraftServerIp == "" {
		log.Printf("[INFO] No minecraft server ip has been provided. Disabling minecraft server status check.")
		return
	}
	if minecraftServerIp != "localhost" && net.ParseIP(minecraftServerIp) == nil {
		log.Printf("[ERROR] Invalid ip \"%s\" for minecraft server has been provided. Disabling minecraft server status check.", minecraftServerIp)
		return
	}

	minecraftServerPort := 25565
	if minecraftServerPortString != "" {
		var err error
		minecraftServerPort, err = strconv.Atoi(minecraftServerPortString)
		if err != nil {
			log.Printf("[ERROR] Non integer port \"%s\" for minecraft server has been provided. Disabling minecraft server status check.", minecraftServerPortString)
			return
		}
	}

	for {
		res, err := minequery.Ping17(minecraftServerIp, minecraftServerPort)
		if err != nil {
			s.UpdateCustomStatus("")
			time.Sleep(QUERY_TIMEOUT)
			continue
		}

		minecraftServerDescription := res.Description.String()
		minecraftServerVersion := res.VersionName
		minecraftServerPlayersMax := res.MaxPlayers
		minecraftServerPlayersOnline := res.OnlinePlayers
		statusString := fmt.Sprintf(
			"⛏️ %s • %s • %d/%d • IP: %s:%s",
			minecraftServerVersion,
			minecraftServerDescription,
			minecraftServerPlayersOnline,
			minecraftServerPlayersMax,
			minecraftServerIp,
			minecraftServerPortString,
		)
		s.UpdateCustomStatus(statusString)
		time.Sleep(QUERY_TIMEOUT)
	}
}
