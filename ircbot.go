package main

import (
	"log"
	"strings"

	"gopkg.in/gcfg.v1"

	irc "github.com/fluffle/goirc/client"
)

// Config is an exported type for defining
// the bot IRC server info
// EXAMPLE:
// [IRC]
// Nickname=bot_nickname
// Host=irc_server
// Channel="#channel_name"
type Config struct {
	IRC struct{ Host, Channel, Nickname string }
}

func main() {
	quit := make(chan bool, 1)

	var config Config
	if err := gcfg.ReadFileInto(&config, "config.gcfg"); err != nil {
		log.Fatal("Couldn't read configuration")
	}

	ircConf := irc.NewConfig(config.IRC.Nickname)
	ircConf.Server = config.IRC.Host
	bot := irc.Client(ircConf)

	bot.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Connected to IRC server", config.IRC.Host)
		conn.Join(config.IRC.Channel)
	})

	bot.HandleFunc("privmsg", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Received:", line.Nick, line.Text())
		if strings.HasPrefix(line.Text(), config.IRC.Nickname) {
			command := strings.Split(line.Text(), " ")[1]
			switch command {
			case "quit":
				log.Println("Received command to quit")
				quit <- true
			}
		}
	})

	log.Println("Connecting to IRC server", config.IRC.Host)
	if err := bot.Connect(); err != nil {
		log.Fatal("IRC connection failed:", err)
	}

	<-quit
}
