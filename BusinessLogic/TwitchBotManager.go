package BusinessLogic

import (
	"TwitchChatBot/Logging"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/TwitchAPI"
	"regexp"
	"time"
)

const ping = "PING :tmi.twitch.tv"

var chatMessageRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG|WHISPER) #*\w+ :!card (.*)?`)

type ITwitchBotManger interface {
	StartTwitchBot(done chan bool)
}

type TwitchBotManager struct {
	TwitchClient TwitchAPI.ITwitchClient
	MagicClient  MagicAPI.IMagicClient
	Logger       Logging.ILogger
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

	// TODO: Get from configuration
	err = this.TwitchClient.Authenticate("magiccardbot", "oauth:b1f8f6mts0e7bxgv4botc6xxtydznz")

	if err != nil {
		this.Logger.Log("Exiting program")
		return
	}

	// TODO: Get channel from configuration
	err = this.TwitchClient.JoinChannel("crendor")

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
			cardName := matches[3]

			go this.lookupCardAndPost(cardName, messageType, user)
		}

		time.Sleep(1 * time.Second)
	}
}

func (this *TwitchBotManager) lookupCardAndPost(cardName string, messageType string, user string) {
	this.Logger.Log("Looking up card " + cardName + " and replying to " + messageType + " and user " + user)
}
