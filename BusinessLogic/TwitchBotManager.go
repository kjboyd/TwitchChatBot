package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"regexp"
)

const ping = "PING :tmi.twitch.tv"

var chatMessageRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG|WHISPER) #*(\w+) :!card (.*)?`)

type ITwitchBotManager interface {
	StartTwitchBot(done chan bool)
}

func NewTwitchBotManager(twitchClient TwitchAPI.ITwitchClient, magicClient MagicAPI.IMagicClient, settings *Configuration.Settings, logger Logging.ILogger) ITwitchBotManager {
	manager := new(twitchBotManager)
	manager.TwitchClient = twitchClient
	manager.MagicClient = magicClient
	manager.Settings = settings
	manager.Logger = logger
	return manager
}

type twitchBotManager struct {
	TwitchClient TwitchAPI.ITwitchClient
	MagicClient  MagicAPI.IMagicClient
	Settings     *Configuration.Settings
	Logger       Logging.ILogger
}

func (this *twitchBotManager) StartTwitchBot(done chan bool) {
	defer func() {
		this.TwitchClient.Disconnect()
		done <- true
	}()

	if err := this.TwitchClient.ConnectToIrcServer(); err != nil {
		this.Logger.Log("Failed to connect to irc channel. Error: " + err.Error())
		return
	}

	if err := this.TwitchClient.Authenticate(this.Settings.UserName, this.Settings.AuthToken); err != nil {
		this.Logger.Log("Failed to authenticate with Twitch. Error: " + err.Error())
		return
	}

	// This is to be able to receive whispers
	if err := this.TwitchClient.JoinChannel(this.Settings.UserName); err != nil {
		this.Logger.Log("Failed to join whispers channel. Error: " + err.Error())
		return
	}

	if err := this.TwitchClient.JoinChannel(this.Settings.Channel); err != nil {
		this.Logger.Log("Failed to join channel. Error: " + err.Error())
		return
	}

	this.monitorChat()
	return
}

func (this *twitchBotManager) monitorChat() {

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
	}
}

func (this *twitchBotManager) lookupCardAndPost(cardName string, messageType string, channel string, user string) {

	this.Logger.Log("Looking up card " + cardName + " and replying to " + messageType + " on channel " + channel + " and user " + user)

	card, err := this.MagicClient.LookupCardInformation(cardName)

	if err != nil {
		go this.TwitchClient.WriteMessage("Unable to find card "+cardName, channel, messageType, user)
		return
	}

	this.Logger.Log("Found card: " + card.String())
	go this.TwitchClient.WriteMessage(card.String(), channel, messageType, user)
}
