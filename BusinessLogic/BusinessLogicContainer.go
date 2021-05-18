package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
)

func NewBusinessLogicContainer(settings *Configuration.Settings) *businessLogicContainer {
	container := new(businessLogicContainer)
	container.TwitchClient = TwitchAPI.NewTwitchClient(settings)
	container.MagicClient = MagicAPI.NewMagicClient(settings)
	container.CardLookupService = NewCardLookupService(container.MagicClient)
	container.ChatBot = NewChatBot(container.TwitchClient, container.CardLookupService, settings)
	return container
}

type businessLogicContainer struct {
	TwitchClient      TwitchAPI.ITwitchClient
	MagicClient       MagicAPI.IMagicClient
	CardLookupService ICardLookupService
	ChatBot           IChatBot
}
