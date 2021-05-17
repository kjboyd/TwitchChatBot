package main

import (
	"TwitchChatBot/BusinessLogic"
	"TwitchChatBot/Configuration"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

var configFilename = flag.String("config", "app.config", "Location of the config file.")

func main() {
	flag.Parse()
	settings, err := ReadConfig(*configFilename)
	if err != nil {
		log.Println("Failed to read config file. Exiting!")
	}

	container := BusinessLogic.NewBusinessLogicContainer(settings)

	done := make(chan bool)
	go BusinessLogic.RunChatBot(container.ChatBot, done)
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
	if err != nil && err != io.EOF {
		return nil, err
	}

	return settings, nil
}
