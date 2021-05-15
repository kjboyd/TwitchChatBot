package MagicAPI

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Infrastructure"
	"TwitchChatBot/Logging"
	"testing"
)

func setupPatient() *magicClient {
	logger := Logging.Logger{}
	settings := Configuration.Settings{
		MagicEndpoint:                 "https://api.magicthegathering.io/v1/",
		MagicRateLimit:                3,
		MagicRateLimitDurationMinutes: 2,
	}
	rateLimiter := Infrastructure.NewRateLimiter(
		settings.MagicRateLimit, settings.MagicRateLimitDurationMinutes*60)
	patient := magicClient{Logger: &logger, Settings: &settings, RateLimiter: rateLimiter}
	return &patient
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

func Test_WillReturnErrorWhenLookingUpNonExistentCard(test *testing.T) {
	patient := setupPatient()

	cardName := "Batman"
	_, err := patient.LookupCardInformation(cardName)

	if err == nil {
		test.Errorf("Error not returned when looking up invalid card.")
		return
	}
}
