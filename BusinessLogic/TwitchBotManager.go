package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"regexp"
	"time"
)

const ping = "PING :tmi.twitch.tv"

var chatMessageRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG|WHISPER) #*(\w+) :!card (.*)?`)

type ITwitchBotManger interface {
	StartTwitchBot(done chan bool)
}

type TwitchBotManager struct {
	TwitchClient TwitchAPI.ITwitchClient
	MagicClient  MagicAPI.IMagicClient
	Logger       Logging.ILogger
	Settings     *Configuration.Settings
}

func (this *TwitchBotManager) StartTwitchBot(done chan bool) {
	defer func() {
		this.TwitchClient.Disconnect()
		done <- true
	}()

	err := this.TwitchClient.ConnectToIrcServer()

	if err != nil {
		this.Logger.Log("Exiting program")
		return
	}

	err = this.TwitchClient.Authenticate(this.Settings.UserName, this.Settings.AuthToken)

	if err != nil {
		this.Logger.Log("Exiting program")
		return
	}

	err = this.TwitchClient.JoinChannel(this.Settings.UserName)
	if err != nil {
		this.Logger.Log("Exiting program")
		return
	}

	err = this.TwitchClient.JoinChannel(this.Settings.Channel)

	if err != nil {
		this.Logger.Log("Exiting program")
		return
	}

	chatChannel := make(chan bool)
	go this.monitorChat(chatChannel)
	<-chatChannel

	return
}

func (this *TwitchBotManager) monitorChat(chatChannel chan bool) {

	defer func() {
		chatChannel <- true
	}()

	for {
		chatLine, err := this.TwitchClient.ReadLine()

		if err != nil {
			this.Logger.Log("Error reading line from chat. " + err.Error())
			return
		}

		if chatLine == ping {
			this.Logger.Log("Sending PONG in response to received PING")
			this.TwitchClient.SendPong()
			continue
		}

		matches := chatMessageRegex.FindStringSubmatch(chatLine)
		if matches != nil {
			user := matches[1]
			messageType := matches[2]
			channel := matches[3]
			cardName := matches[4]

			go this.lookupCardAndPost(cardName, messageType, channel, user)
		}

		this.Logger.Log("Message: " + chatLine)
		time.Sleep(1 * time.Second)
	}
}

func (this *TwitchBotManager) lookupCardAndPost(cardName string, messageType string, channel string, user string) {

	this.Logger.Log("Looking up card " + cardName + " and replying to " + messageType + " on channel " + channel + " and user " + user)

	card, err := this.MagicClient.LookupCardInformation(cardName)

	if err != nil {
		go this.TwitchClient.WriteMessage("Unable to find card "+cardName, channel, messageType, user)
		return
	}
	go this.TwitchClient.WriteMessage(card.String(), channel, messageType, user)
}
