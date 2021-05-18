package BusinessLogic

import (
	"TwitchChatBot/MagicAPI"
	"io"
	"log"
)

type ICardLookupService interface {
	LookupCardAndPost(cardName string, writer io.StringWriter)
}

func NewCardLookupService(magicClient MagicAPI.IMagicClient) ICardLookupService {
	service := new(cardLookupService)
	service.MagicClient = magicClient
	return service
}

type cardLookupService struct {
	MagicClient MagicAPI.IMagicClient
}

func (this *cardLookupService) LookupCardAndPost(
	cardName string, writer io.StringWriter) {

	if cardName == "" {
		go writer.WriteString("Please specify card name.")
		return
	}

	card, err := this.MagicClient.LookupCardInformation(cardName)

	if err != nil {
		go writer.WriteString("Unable to find card " + cardName)
		return
	}

	log.Println("Found card: " + card.String())
	go writer.WriteString(card.String())
}
