package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/TwitchAPI"
	"log"
	"regexp"
)

const ping = "PING :tmi.twitch.tv"

var chatMessageRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG|WHISPER) #*(\w+) :!(\w+) *(.*)?`)

type IChatBot interface {
	Connect() error
	Disconnect()
	ProcessMessage() bool
}

func NewChatBot(twitchClient TwitchAPI.ITwitchClient,
	cardLookupService ICardLookupService,
	settings *Configuration.Settings) IChatBot {
	bot := new(chatBot)
	bot.TwitchClient = twitchClient
	bot.CardLookupService = cardLookupService
	bot.Settings = settings
	bot.CurrentChannel = settings.Channel
	return bot
}

type chatBot struct {
	TwitchClient      TwitchAPI.ITwitchClient
	CardLookupService ICardLookupService
	Settings          *Configuration.Settings
	CurrentChannel    string
}

func (this *chatBot) Connect() error {
	if err := this.TwitchClient.ConnectToIrcServer(); err != nil {
		log.Println("Failed to connect to irc channel. Error: " + err.Error())
		return err
	}

	if err := this.TwitchClient.Authenticate(this.Settings.UserName, this.Settings.AuthToken); err != nil {
		log.Println("Failed to authenticate with Twitch. Error: " + err.Error())
		return err
	}

	// This is to be able to receive whispers
	if err := this.TwitchClient.JoinChannel(this.Settings.UserName); err != nil {
		log.Println("Failed to join whispers channel. Error: " + err.Error())
		return err
	}

	if err := this.TwitchClient.JoinChannel(this.Settings.Channel); err != nil {
		log.Println("Failed to join channel. Error: " + err.Error())
		return err
	}

	return nil
}

func (this *chatBot) Disconnect() {
	this.TwitchClient.LeaveChannel(this.CurrentChannel)
	this.TwitchClient.LeaveChannel(this.Settings.UserName)
	this.TwitchClient.Disconnect()
}

func (this *chatBot) ProcessMessage() bool {

	chatLine, err := this.TwitchClient.ReadLine()

	if err != nil {
		log.Println("Error reading line from chat. " + err.Error())
		return false
	}

	if chatLine == ping {
		log.Println("Sending PONG in response to received PING")
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
			go this.CardLookupService.LookupCardAndPost(
				argument, messageType, channel, user)
		case this.Settings.DisconnectCommand:
			if messageType == TwitchAPI.WhisperMessageType {
				log.Println("Received command to shutdown!")
				return false
			}
		case this.Settings.ChangeChannelCommand:
			if messageType == TwitchAPI.WhisperMessageType {
				// If we leave our own channel, we won't receive whispers anymore
				if this.CurrentChannel != this.Settings.UserName {
					this.TwitchClient.LeaveChannel(this.CurrentChannel)
				}
				this.TwitchClient.JoinChannel(argument)
				this.CurrentChannel = argument
			}
		}
	}
	//log.Println("Message: " + chatLine)

	return true
}
