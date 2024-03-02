package events

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func VoiceStateUpdateHandler(s *discord.Session, e *discord.VoiceStateUpdate) {
	// Return if handling client's own voice state update
	if e.UserID == os.Getenv("CLIENT_ID") {
		return
	}

	err := greetingHandler(s, e)
	if err != nil {
		log.Printf("[ERROR] Could not greet a user: %v", err)
	}

	err = birthdayHandler(s, e)
	if err != nil {
		log.Printf("[ERROR] Could not congratulate a user with his birthday: %v", err)
	}
}

func greetingHandler(s *discord.Session, e *discord.VoiceStateUpdate) error {
	// Return if greeting functionality is disabled
	if os.Getenv("GREETING_TIME_PERIOD") == "-1" {
		return nil
	}

	// Check if user is in a voice channel now
	// Return if user is not connected to a voice channel
	if e.ChannelID == "" {
		return nil
	}

	// Check if user has moved from no voice channel to a voice channel, otherwise return
	if !(e.ChannelID != "" && (e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == "")) {
		return nil
	}

	// Check for supplied database filepath or use the default one
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata.sqlite3.db"
	}

	// Open db connection and create the greeting table if if does not exist
	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	db.AutoMigrate(&utl.GreetedUser{})

	// Fetch user by his userId
	var user utl.GreetedUser
	res := db.First(&user, "discord_user_id = ?", e.UserID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Insert a user
		user = utl.GreetedUser{
			DiscordUserId:         e.UserID,
			GreetingUnixTimestamp: time.Now().Unix(),
		}
		res = db.Create(&user)
		if res.Error != nil {
			return fmt.Errorf("error inserting a new user into database: %v", res.Error)
		}

		// Greet a user
		err = greet(s, e)
		if err != nil {
			return fmt.Errorf("error while greeting: %v", err)
		}
	}

	lastGreetingUnixTime := user.GreetingUnixTimestamp

	// Get greeting time petiod env. variable or use default (one week)
	var greetingTimePeriod int64 = 604800 // One week in unix seconds
	if os.Getenv("GREETING_TIME_PERIOD") != "" {
		greetingTimePeriod, err = strconv.ParseInt(os.Getenv("GREETING_TIME_PERIOD"), 10, 64)
		if err != nil {
			greetingTimePeriod = 604800 // One week in unix seconds
			log.Printf("[ERROR] Error parsing GREETING_TIME_PERIOD env. variable; using default (one week): %v", err)
		}
	}

	// Check if the time has come to greet user again
	if time.Now().Unix()-lastGreetingUnixTime > greetingTimePeriod {
		// Update user's greeting time
		err = updateUserGreetingTime(db, e.UserID, time.Now().Unix())
		if err != nil {
			return fmt.Errorf("error updating user's greeting time: %v", err)
		}

		// Greet the user
		err = greet(s, e)
		if err != nil {
			return fmt.Errorf("error reading rows: %v", err)
		}
	}

	return nil
}

func updateUserGreetingTime(db *gorm.DB, discordUserId string, unixTimestamp int64) error {
	var user utl.GreetedUser
	res := db.First(&user, "discord_user_id = ?", discordUserId)
	if res.Error != nil {
		return fmt.Errorf("error updating user in database: %v", res.Error)
	}

	user.GreetingUnixTimestamp = unixTimestamp
	res = db.Save(&user)
	if res.Error != nil {
		return fmt.Errorf("error updating user in database: %v", res.Error)
	}

	return nil
}

func greet(s *discord.Session, e *discord.VoiceStateUpdate) error {
	if e.ChannelID == "" {
		return nil
	}

	const GREETING_FILEPATH = "./assets/audio/greeting.opus"

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
	// Check if user is in a voice channel now
	// Return if user is not connected to a voice channel
	if e.ChannelID == "" {
		return nil
	}

	// Check if user has moved from no voice channel to a voice channel, otherwise return
	if !(e.ChannelID != "" && (e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == "")) {
		return nil
	}

	// Check for supplied database filepath or use the default one
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata.sqlite3.db"
	}

	// Open db connection and create the birthday table if if does not exist
	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	db.AutoMigrate(&utl.BirthdayUser{})

	// Fetch user by his userId
	var user utl.BirthdayUser
	res := db.First(&user, "discord_user_id = ?", e.UserID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// If user has not been found, return
		log.Println("[INFO] No birthday user has been found in the datbase.")
		return nil
	}

	// Congratulate user if it is his birthday
	if user.BirthdayDate == time.Now().Format("01/02") && user.LastGreetingYear < time.Now().Year() {
		// Set user's last greeted year to this year
		user.LastGreetingYear = time.Now().Year()
		res = db.Save(&user)
		if res.Error != nil {
			return fmt.Errorf("error updating last greeting year in database: %v", res.Error)
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

	// Pick random file from directory
	files, err := os.ReadDir(BIRTHDAY_DIR_PATH)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	filename := BIRTHDAY_DIR_PATH + filenames[rand.Intn(len(filenames))]

	// Read and play audio file
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
