package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
)

func NewBusinessLogicContainer(settings *Configuration.Settings, logger Logging.ILogger) *businessLogicContainer {
	container := new(businessLogicContainer)
	container.TwitchClient = TwitchAPI.NewTwitchClient(settings, logger)
	container.MagicClient = MagicAPI.NewMagicClient(settings, logger)
	container.ChatBot = NewChatBot(container.TwitchClient, container.MagicClient, settings, logger)
	return container
}

type businessLogicContainer struct {
	TwitchClient TwitchAPI.ITwitchClient
	MagicClient  MagicAPI.IMagicClient
	ChatBot      IChatBot
}
