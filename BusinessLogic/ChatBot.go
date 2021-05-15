package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"regexp"
)

const ping = "PING :tmi.twitch.tv"

var chatMessageRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG|WHISPER) #*(\w+) :!(\w+) *(.*)?`)

type IChatBot interface {
	Connect() error
	Disconnect()
	ProcessMessage() bool
}

func NewChatBot(twitchClient TwitchAPI.ITwitchClient, magicClient MagicAPI.IMagicClient, settings *Configuration.Settings, logger Logging.ILogger) *chatBot {
	bot := new(chatBot)
	bot.TwitchClient = twitchClient
	bot.MagicClient = magicClient
	bot.Settings = settings
	bot.Logger = logger
	bot.CurrentChannel = settings.Channel
	return bot
}

type chatBot struct {
	TwitchClient   TwitchAPI.ITwitchClient
	MagicClient    MagicAPI.IMagicClient
	Settings       *Configuration.Settings
	Logger         Logging.ILogger
	CurrentChannel string
}

func (this *chatBot) Connect() error {
	if err := this.TwitchClient.ConnectToIrcServer(); err != nil {
		this.Logger.Log("Failed to connect to irc channel. Error: " + err.Error())
		return err
	}

	if err := this.TwitchClient.Authenticate(this.Settings.UserName, this.Settings.AuthToken); err != nil {
		this.Logger.Log("Failed to authenticate with Twitch. Error: " + err.Error())
		return err
	}

	// This is to be able to receive whispers
	if err := this.TwitchClient.JoinChannel(this.Settings.UserName); err != nil {
		this.Logger.Log("Failed to join whispers channel. Error: " + err.Error())
		return err
	}

	if err := this.TwitchClient.JoinChannel(this.Settings.Channel); err != nil {
		this.Logger.Log("Failed to join channel. Error: " + err.Error())
		return err
	}

	return nil
}

func (this *chatBot) Disconnect() {
	this.TwitchClient.LeaveChannel(this.CurrentChannel)
	this.TwitchClient.Disconnect()
}

func (this *chatBot) ProcessMessage() bool {

	chatLine, err := this.TwitchClient.ReadLine()

	if err != nil {
		this.Logger.Log("Error reading line from chat. " + err.Error())
		return false
	}

	if chatLine == ping {
		this.Logger.Log("Sending PONG in response to received PING")
		this.TwitchClient.SendPong()
		return true
	}

	matches := chatMessageRegex.FindStringSubmatch(chatLine)
	if matches != nil {
		user := matches[1]
		messageType := matches[2]
		channel := matches[3]
		command := matches[4]
		argument := ""
		if len(matches) > 5 {
			argument = matches[5]
		}

		switch command {
		case this.Settings.CardCommand:
			go this.lookupCardAndPost(argument, messageType, channel, user)
		case this.Settings.DisconnectCommand:
			this.Logger.Log("Received command to shutdown!")
			return false
		case this.Settings.ChangeChannelCommand:
			this.TwitchClient.LeaveChannel(this.CurrentChannel)
			this.TwitchClient.JoinChannel(argument)
			this.CurrentChannel = argument
		}
	}
	//this.Logger.Log("Message: " + chatLine)

	return true
}

func (this *chatBot) lookupCardAndPost(cardName string, messageType string, channel string, user string) {

	this.Logger.Log("Looking up card " + cardName + " and replying to " + messageType + " on channel " + channel + " and user " + user)

	if cardName == "" {
		go this.TwitchClient.WriteMessage("Please specify card name.", channel, messageType, user)
		return
	}

	card, err := this.MagicClient.LookupCardInformation(cardName)

	if err != nil {
		go this.TwitchClient.WriteMessage("Unable to find card "+cardName, channel, messageType, user)
		return
	}

	this.Logger.Log("Found card: " + card.String())
	go this.TwitchClient.WriteMessage(card.String(), channel, messageType, user)
}
