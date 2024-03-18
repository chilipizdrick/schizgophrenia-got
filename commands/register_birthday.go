package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var RegisterBirthdayCommand = utl.SlashCommand{
	CommandData: &discord.ApplicationCommand{
		Name:        "register-birthday",
		Description: "Register discord user's birthday.",
		Options: []*discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionString,
				Name:        "user",
				Description: "Discord user mention",
				Required:    true,
			},
			{
				Type:        discord.ApplicationCommandOptionString,
				Name:        "birthday",
				Description: "User's birthday date in MM/DD format",
				Required:    true,
			},
		},
	},
	CommandHandler: func(s *discord.Session, i *discord.InteractionCreate) {
		optionMap := utl.GetOptionMap(i)
		userMention := optionMap["user"].StringValue()

		// Check if provided field is a discord user mention
		if len(userMention) != 21 || !(strings.HasPrefix(userMention, "<@") && strings.HasSuffix(userMention, ">")) {
			log.Println("[ERROR] Supplied \"user\" field is not a discord user mention.")
			utl.RespondToInteractionCreateWithString(s, i, "Provided \"user\" field is not a valid discord user mention.")
			return
		}

		// Parse discord user id from user mention
		userID := userMention[2 : len(userMention)-1]

		birthdayDate := optionMap["birthday"].StringValue()

		// Validate provided date
		_, err := time.Parse("01/02", birthdayDate)
		if err != nil {
			log.Println("[ERROR] Supplied \"birthday\" field is not a valid date.")
			utl.RespondToInteractionCreateWithString(s, i, "Provided \"birthday\" field is not a valid date.")
			return
		}

		// Check for supplied database filepath or use the default one
		sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
		if sqliteDatabaseFilepath == "" {
			sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
		}

		// Open db connection and create the birthday table if if does not exist
		db, err := gorm.Open(sqlite.Open(sqliteDatabaseFilepath), &gorm.Config{})
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "An error occured while executing command.")
			log.Printf("error opening database connection: %v", err)
			return
		}
		db.AutoMigrate(&utl.BirthdayUser{})

		// Fetch user by his userId
		var user utl.BirthdayUser
		res := db.First(&user, "discord_user_id = ?", userID)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// If user has not been found, add it to batabase
			user = utl.BirthdayUser{
				DiscordUserId: userID,
				BirthdayDate:  birthdayDate,
			}
			res = db.Create(&user)
			if res.Error != nil {
				utl.RespondToInteractionCreateWithString(s, i, "An error occured while executing command.")
				log.Printf("[ERROR] Could not create user %v", res.Error)
				return
			}
			utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated <@%v>'s birthday to be %s.", userID, birthdayDate))
			return
		}

		// Otherwise update user's greeting date
		user.BirthdayDate = birthdayDate
		res = db.Save(&user)
		if res.Error != nil {
			utl.RespondToInteractionCreateWithString(s, i, "An error occured while executing command.")
			log.Printf("[ERROR] Could not update user's birthday date %v", res.Error)
			return
		}

		utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated <@%v>'s birthday to be %s.", userID, birthdayDate))
	},
}
