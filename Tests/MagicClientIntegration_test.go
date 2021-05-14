package Tests

import (
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"testing"
)

func setupPatient() MagicAPI.MagicClient {
	Logger := Logging.Logger{}
	patient := MagicAPI.MagicClient{
		Logger:        &Logger,
		MagicEndpoint: "https://api.magicthegathering.io/v1/",
	}
	return patient
}

func Test_WillLookupCardInformationById(test *testing.T) {
	patient := setupPatient()

	card, err := patient.LookupCardInformation("409741")

	if err != nil {
		test.Errorf("Error when looking up card information by Id.")
	}

	if card.Name != "Archangel Avacyn // Avacyn, the Purifier" {
		test.Errorf("Card information is not correct.")
	}
}

func Test_WillLookupCardInformationByName(test *testing.T) {
	patient := setupPatient()

	cardName := "Academic Probation"
	card, err := patient.LookupCardInformation(cardName)

	if err != nil {
		test.Errorf("Error when looking up card information by name.")
		return
	}

	if card.Name != cardName {
		test.Errorf("Card information is not correct. Expected " + cardName + " Found " + card.String())
		return
	}
}
