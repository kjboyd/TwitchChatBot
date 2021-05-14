package main

import (
	"TwitchChatBot/BusinessLogic"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
)

func main() {
	Logger := Logging.Logger{}
	twitchClient := TwitchAPI.TwitchClient{
		Port:   "6667",
		Server: "irc.chat.twitch.tv",
		Logger: &Logger,
	}
	magicClient := MagicAPI.MagicClient{}

	twitchBotManager := BusinessLogic.TwitchBotManager{
		MagicClient:  &magicClient,
		TwitchClient: &twitchClient,
		Logger:       &Logger,
	}
	done := make(chan bool)
	go twitchBotManager.StartTwitchBot(done)
	<-done
}
