package utils

type BirthdayUser struct {
	ID               uint
	DiscordUserId    string
	BirthdayDate     string // In MM/DD format
	LastGreetingYear int
}

type GreetedUser struct {
	ID                    uint
	DiscordUserId         string
	GreetingUnixTimestamp int64
}
