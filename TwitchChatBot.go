package main

import (
	"TwitchChatBot/BusinessLogic"
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

func main() {
	logger := Logging.Logger{}
	settings, err := ReadConfig("app.config")
	if err != nil {
		logger.Log("Failed to read config file. Exiting!")
	}

	twitchClient := TwitchAPI.TwitchClient{
		Port:   settings.TwitchPort,
		Server: settings.TwitchServer,
		Logger: &logger,
	}
	magicClient := MagicAPI.MagicClient{
		Logger:        &logger,
		MagicEndpoint: settings.MagicEndpoint,
	}

	twitchBotManager := BusinessLogic.TwitchBotManager{
		MagicClient:  &magicClient,
		TwitchClient: &twitchClient,
		Logger:       &logger,
		Settings:     settings,
	}
	done := make(chan bool)
	go twitchBotManager.StartTwitchBot(done)
	<-done
}

func ReadConfig(configPath string) (*Configuration.Settings, error) {

	// reads from the file
	settingsFile, err := ioutil.ReadFile(configPath)
	if nil != err {
		return nil, err
	}

	settings := &Configuration.Settings{}

	// parses the file contents
	dec := json.NewDecoder(strings.NewReader(string(settingsFile)))
	err = dec.Decode(settings)
	if nil != err && io.EOF != err {
		return nil, err
	}

	return settings, nil
}
