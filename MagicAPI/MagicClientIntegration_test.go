package MagicAPI

import (
	"TwitchChatBot/Configuration"
	"testing"
)

func setupPatient() IMagicClient {
	settings := Configuration.Settings{
		MagicEndpoint:                 "https://api.magicthegathering.io/v1/",
		MagicRateLimit:                3,
		MagicRateLimitDurationSeconds: 2,
	}
	patient := NewMagicClient(&settings)
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

func Test_WillReturnErrorWhenLookingUpNonExistentCardName(test *testing.T) {
	patient := setupPatient()

	cardName := "Batman"
	_, err := patient.LookupCardInformation(cardName)

	if err == nil {
		test.Errorf("Error not returned when looking up invalid card name.")
		return
	}
}

func Test_WillReturnErrorWhenLookingUpNonExistentMultiverseId(test *testing.T) {
	patient := setupPatient()

	multiverseId := "0"
	_, err := patient.LookupCardInformation(multiverseId)

	if err == nil {
		test.Errorf("Error not returned when looking up invalid multiverse Id.")
		return
	}
}
