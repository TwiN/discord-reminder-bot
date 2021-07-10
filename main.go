package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TwinProduction/discord-reminder-bot/config"
	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/TwinProduction/discord-reminder-bot/discord"
	"github.com/bwmarrin/discordgo"
)

var (
	killChannel chan os.Signal
	cfg         *config.Config
)

func main() {
	err := database.Initialize("sqlite", "data.db")
	if err != nil {
		panic(err)
	}
	cfg = config.Get()
	bot, err := Connect(cfg.DiscordToken)
	if err != nil {
		panic(err)
	}
	defer bot.Close()
	log.Printf("Bot with id=%s has connected successfully", bot.State.User.ID)
	discord.Start(bot, cfg)
	waitUntilTermination()
	log.Println("Terminating bot")
}

func waitUntilTermination() {
	killChannel = make(chan os.Signal, 1)
	signal.Notify(killChannel, syscall.SIGTERM)
	<-killChannel
}

// Connect starts a Discord session
func Connect(discordToken string) (*discordgo.Session, error) {
	discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions)
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}
	err = session.Open()
	return session, err
}
