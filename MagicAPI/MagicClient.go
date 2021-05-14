package MagicAPI

import (
	"TwitchChatBot/Logging"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type IMagicClient interface {
	LookupCardInformation(cardNameOrId string) (*MagicCard, error)
}

type MagicClient struct {
	Logger        Logging.ILogger
	MagicEndpoint string
}

func (this *MagicClient) LookupCardInformation(cardNameOrId string) (*MagicCard, error) {

	cardId, err := strconv.Atoi(cardNameOrId)

	if err == nil {
		return this.lookupCardById(fmt.Sprint(cardId))
	} else {
		return this.lookupCardByName(cardNameOrId)
	}
}

func (this *MagicClient) lookupCardByName(cardName string) (*MagicCard, error) {
	this.Logger.Log("Looking up card by name: " + cardName)

	resp, err := http.Get(this.MagicEndpoint + "cards?name=" + url.QueryEscape(cardName))
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

func (this *MagicClient) lookupCardById(cardId string) (*MagicCard, error) {
	this.Logger.Log("Looking up card by id: " + cardId)

	resp, err := http.Get(this.MagicEndpoint + "cards/" + cardId)
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

func (this *MagicClient) decodeCard(body io.Reader) (*MagicCard, error) {
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

	if cardResponse.Card != nil {
		return cardResponse.Card, nil
	}
	if len(cardResponse.Cards) > 0 {
		return cardResponse.Cards[0], nil
	}
	return nil, fmt.Errorf("Card not found")
}
