package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TwinProduction/discord-reminder-bot/config"
	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

var (
	killChannel chan os.Signal
	botMention  string
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
	botMention = "<@!" + bot.State.User.ID + ">"
	bot.AddHandler(HandleMessage)
	bot.AddHandler(HandleReactionAdd)
	_ = bot.UpdateListeningStatus(config.Get().CommandPrefix + "RemindMe")
	log.Printf("Bot with id=%s has connected successfully", bot.State.User.ID)
	go worker(bot)
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
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}
	err = discord.Open()
	return discord, err
}
