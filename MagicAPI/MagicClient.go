package MagicAPI

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Infrastructure"
	"TwitchChatBot/Logging"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type IMagicClient interface {
	LookupCardInformation(cardNameOrId string) (*MagicCard, error)
}

func NewMagicClient(settings *Configuration.Settings, logger Logging.ILogger) IMagicClient {
	client := new(magicClient)
	client.Logger = logger
	client.Settings = settings
	client.RateLimiter = Infrastructure.NewRateLimiter(
		settings.MagicRateLimit, time.Duration(settings.MagicRateLimitDurationSeconds)*time.Second)
	return client
}

type magicClient struct {
	Logger      Logging.ILogger
	Settings    *Configuration.Settings
	RateLimiter Infrastructure.IRateLimiter
}

func (this *magicClient) LookupCardInformation(cardNameOrId string) (*MagicCard, error) {

	var result *MagicCard = nil
	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		cardId, err := strconv.Atoi(cardNameOrId)

		if err == nil {
			result, err = this.lookupCardById(fmt.Sprint(cardId))
		} else {
			result, err = this.lookupCardByName(cardNameOrId)
		}
	})
	return result, err
}

func (this *magicClient) lookupCardByName(cardName string) (*MagicCard, error) {

	this.Logger.Log("Looking up card by name: " + cardName)

	resp, err := http.Get(this.Settings.MagicEndpoint + "cards?name=" + url.QueryEscape(cardName))
	if err != nil {
		this.Logger.Log("Got error when looking up card with name " + cardName + ". Error: " + err.Error())
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	if resp.StatusCode != 200 {
		this.Logger.Log("Got response " + fmt.Sprint(resp.StatusCode) + " when looking up card with name " + cardName)
		return nil, fmt.Errorf("Card not found")
	}

	card, err := this.decodeCard(body)
	if err != nil {
		this.Logger.Log("Unable to lookup card with name: " + cardName)
		return nil, err
	}

	return card, nil
}

func (this *magicClient) lookupCardById(cardId string) (*MagicCard, error) {
	this.Logger.Log("Looking up card by id: " + cardId)

	resp, err := http.Get(this.Settings.MagicEndpoint + "cards/" + cardId)
	if err != nil {
		this.Logger.Log("Got error when looking up card with Id " + cardId + ". Error: " + err.Error())
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	if resp.StatusCode != 200 {
		this.Logger.Log("Got response " + fmt.Sprint(resp.StatusCode) + " when looking up card with Id " + cardId)
		return nil, fmt.Errorf("Card not found")
	}

	card, err := this.decodeCard(body)
	if err != nil {
		this.Logger.Log("Unable to lookup card with Id: " + cardId)
		return nil, err
	}

	return card, nil
}

func (this *magicClient) decodeCard(body io.Reader) (*MagicCard, error) {
	cardResponse := new(CardResponse)
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&cardResponse)

	if err != nil {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			this.Logger.Log("Error decoding card")
			return nil, err
		}
		bodyString := string(bodyBytes)
		this.Logger.Log("Error decoding card from: " + bodyString)
		return nil, err
	}

	// The Magic API can return multiple values for a query, but we are always querying on a single card name
	if cardResponse.Card != nil {
		return cardResponse.Card, nil
	}
	if len(cardResponse.Cards) > 0 {
		return cardResponse.Cards[0], nil
	}
	return nil, fmt.Errorf("Card not found")
}
