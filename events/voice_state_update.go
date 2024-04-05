package events

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
)

func VoiceStateUpdateHandler(s *discord.Session, e *discord.VoiceStateUpdate) {
	// Return if handling client's own voice state update
	if e.UserID == os.Getenv("CLIENT_ID") {
		return
	}

	guildData, err := utl.LoadGuildFromDBByID(e.GuildID)
	if err != nil {
		log.Printf("[ERROR] Could not load guild from database: %e", err)
	}

	if guildData.Greeting {
		err = greetingHandler(s, e)
		if err != nil {
			log.Printf("[ERROR] Could not greet a user: %v", err)
		}
	}

	if guildData.Birthbay {
		err = birthdayHandler(s, e)
		if err != nil {
			log.Printf("[ERROR] Could not congratulate a user with his birthday: %v", err)
		}
	}
}

func greetingHandler(s *discord.Session, e *discord.VoiceStateUpdate) error {
	// Return if greeting functionality is disabled
	if os.Getenv("GREETING_TIME_PERIOD") == "-1" {
		return nil
	}

	// Return if user is not connected to a voice channel
	if e.ChannelID == "" {
		return nil
	}

	// Check if user has moved from no voice channel to a voice channel, otherwise return
	if !(e.ChannelID != "" && (e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == "")) {
		return nil
	}

	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	userData, err := utl.LoadUserFromDBByID(e.UserID)
	if err != nil {
		return fmt.Errorf("error loading user from database: %v", err)
	}

	lastGreetingUnixTime := userData.GreetingUnixTimestamp

	var greetingTimePeriod int64 = 604800 // One week in unix seconds
	if os.Getenv("GREETING_TIME_PERIOD") != "" {
		greetingTimePeriod, err = strconv.ParseInt(os.Getenv("GREETING_TIME_PERIOD"), 10, 64)
		if err != nil {
			greetingTimePeriod = 604800 // One week in unix seconds
			log.Printf("[ERROR] Error parsing GREETING_TIME_PERIOD env. variable; using default (one week): %v", err)
		}
	}

	if time.Now().Unix()-lastGreetingUnixTime > greetingTimePeriod {
		userData.GreetingUnixTimestamp = time.Now().Unix()
		err = utl.SaveUserToDB(userData)
		if err != nil {
			return fmt.Errorf("error updating user's greeting time: %v", err)
		}

		err = greet(s, e)
		if err != nil {
			return fmt.Errorf("error reading rows: %v", err)
		}
	}

	return nil
}

func greet(s *discord.Session, e *discord.VoiceStateUpdate) error {
	if e.ChannelID == "" {
		return nil
	}

	const GREETING_FILEPATH = "./assets/audio/greeting.ogg"

	var audioBuffer [][]byte
	err := utl.LoadOpusFile(GREETING_FILEPATH, &audioBuffer)
	if err != nil {
		return fmt.Errorf("error loading opus file: %v", err)
	}

	err = utl.PlayAudio(s, e.GuildID, e.ChannelID, audioBuffer)
	if err != nil {
		return fmt.Errorf("error playing audio: %v", err)
	}

	return nil
}

func birthdayHandler(s *discord.Session, e *discord.VoiceStateUpdate) error {
	// Return if user is not connected to a voice channel
	if e.ChannelID == "" {
		return nil
	}

	// Check if user has moved from no voice channel to a voice channel, otherwise return
	if !(e.ChannelID != "" && (e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == "")) {
		return nil
	}

	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	userData, err := utl.LoadUserFromDBByID(e.UserID)
	if err != nil {
		return fmt.Errorf("error loading user from databse: %v", err)
	}

	if userData.BirthdayDate == time.Now().Format("01/02") &&
		userData.LastBirthdayGreetingYear < time.Now().Year() {

		userData.LastBirthdayGreetingYear = time.Now().Year()
		err = utl.SaveUserToDB(userData)
		if err != nil {
			return fmt.Errorf("error updating last birthday greeting year in database: %v", err)
		}

		err = congratulate(s, e)
		if err != nil {
			return fmt.Errorf("error while congratulating: %v", err)
		}
	}

	return nil
}

func congratulate(s *discord.Session, e *discord.VoiceStateUpdate) error {
	if e.ChannelID == "" {
		return nil
	}

	const BIRTHDAY_DIR_PATH = "./assets/audio/birthday/"

	filename, err := utl.PickRandomFileFromDirectory(BIRTHDAY_DIR_PATH)
	if err != nil {
		return fmt.Errorf("error picking random file from directory: %v", err)
	}

	var audioBuffer [][]byte
	err = utl.LoadOpusFile(filename, &audioBuffer)
	if err != nil {
		return fmt.Errorf("error loading opus file: %v", err)
	}

	err = utl.PlayAudio(s, e.GuildID, e.ChannelID, audioBuffer)
	if err != nil {
		return fmt.Errorf("error playing audio: %v", err)
	}

	return nil
}
