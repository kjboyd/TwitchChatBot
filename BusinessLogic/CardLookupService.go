package BusinessLogic

import (
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"log"
)

type ICardLookupService interface {
	LookupCardAndPost(cardName string, messageType string, channel string, user string)
}

func NewCardLookupService(twitchClient TwitchAPI.ITwitchClient, magicClient MagicAPI.IMagicClient) ICardLookupService {
	service := new(cardLookupService)
	service.TwitchClient = twitchClient
	service.MagicClient = magicClient
	return service
}

type cardLookupService struct {
	TwitchClient TwitchAPI.ITwitchClient
	MagicClient  MagicAPI.IMagicClient
}

func (this *cardLookupService) LookupCardAndPost(
	cardName string, messageType string, channel string, user string) {

	log.Println("Looking up card " + cardName + " and replying to " + messageType + " on channel " + channel + " for user " + user)

	if cardName == "" {
		go this.TwitchClient.WriteMessage("Please specify card name.", channel, messageType, user)
		return
	}

	card, err := this.MagicClient.LookupCardInformation(cardName)

	if err != nil {
		go this.TwitchClient.WriteMessage("Unable to find card "+cardName, channel, messageType, user)
		return
	}

	log.Println("Found card: " + card.String())
	go this.TwitchClient.WriteMessage(card.String(), channel, messageType, user)
}
