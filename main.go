package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discordgo "github.com/bwmarrin/discordgo"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

var config Config

func main() {
	fmt.Println("Hello World")
	//Open Config
	configFile, err := os.Open("config.json")
	//Check for error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully opened config file")
	//defer closing of file
	defer configFile.Close()

	//Read in config
	byteValue, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(byteValue, &config)
	//For testing reading of token
	fmt.Println(config.Token)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Ignore messages created by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Commands
	if strings.HasPrefix(m.Content, config.Prefix) {
		m.Content = strings.TrimPrefix(m.Content, config.Prefix)
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
		if strings.HasPrefix(m.Content, "MD5") {
			preHash := strings.TrimPrefix(m.Content, "MD5 ")
			hash := md5.Sum([]byte(preHash))
			s.ChannelMessageSend(m.ChannelID, hex.EncodeToString(hash[:]))
		}
	}
}
