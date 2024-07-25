package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	discord "github.com/bwmarrin/discordgo"
	utl "github.com/chilipizdrick/schizgophrenia-got/utils"
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

		sqliteDatabaseFilepath := os.Getenv("SQLITE_DATABASE_FILEPATH")
		if sqliteDatabaseFilepath == "" {
			sqliteDatabaseFilepath = "./userdata/userdata.sqlite3.db"
		}

		userData, err := utl.LoadUserFromDBByID(userID)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "An error occured while executing command.")
			log.Printf("[ERROR] Could not load user from database: %s", err)
			return
		}

		userData.BirthdayDate = birthdayDate
		err = utl.SaveUserToDB(userData)
		if err != nil {
			utl.RespondToInteractionCreateWithString(s, i, "An error occured while executing command.")
			log.Printf("[ERROR] Could not update user's birthday date: %s", err)
			return
		}

		utl.RespondToInteractionCreateWithString(s, i, fmt.Sprintf("Successfully updated <@%s>'s birthday to be %s.", userID, birthdayDate))
	},
}
