package utils

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID                       uint
	DiscordUserID            string
	DiscordGuildID           string
	GreetingEnabled          bool
	GreetingUnixTimestamp    int64
	BirthdayDate             string // In MM/DD format
	LastBirthdayGreetingYear int
}

type Guild struct {
	ID             uint
	DiscordGuildID string
	Greeting       bool
	Birthday       bool
	DMOnMention    bool
}

func LoadUserFromDBByID(userID string) (*User, error) {
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error opening database: %e", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("error migrating \"users\" schema: %e", err)
	}

	var userData User
	res := db.First(&userData, "discord_user_id = ?", userID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Printf("[INFO] Could not find user (discord_user_id=%s) in the database, creating new.", userID)
		userData = User{
			DiscordUserID:   userID,
			GreetingEnabled: true,
		}

		res = db.Create(&userData)
		if res.Error != nil {
			return nil, fmt.Errorf("error inserting a new user into database: %e", res.Error)
		}
	}

	return &userData, nil
}

func SaveUserToDB(userData *User) error {
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database: %e", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("error migrating \"user\" schema: %e", err)
	}

	res := db.Save(userData)
	if res.Error != nil {
		return fmt.Errorf("error saving user to database: %e", res.Error)
	}

	return nil
}

func LoadGuildFromDBByID(guildID string) (*Guild, error) {
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error opening database: %e", err)
	}

	err = db.AutoMigrate(&Guild{})
	if err != nil {
		return nil, fmt.Errorf("error migrating \"guilds\" schema: %e", err)
	}

	var guildData Guild
	res := db.First(&guildData, "discord_guild_id = ?", guildID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Printf("[INFO] Could not find guild (discord_guild_id=%s) in the database, creating new.", guildID)
		guildData = Guild{
			DiscordGuildID: guildID,
			Greeting:       false,
			Birthday:       false,
		}

		res = db.Create(&guildData)
		if res.Error != nil {
			return nil, fmt.Errorf("error inserting a new guild into database: %e", res.Error)
		}
	}

	return &guildData, nil
}

func SaveGuildToDB(guildData *Guild) error {
	sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
	if sqliteDatabaseFilepath == "" {
		sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
	}

	db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database: %e", err)
	}

	err = db.AutoMigrate(&Guild{})
	if err != nil {
		return fmt.Errorf("error migrating \"guilds\" schema: %e", err)
	}

	res := db.Save(guildData)
	if res.Error != nil {
		return fmt.Errorf("error saving guild to database: %e", res.Error)
	}

	return nil
}
