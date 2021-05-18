package TwitchAPI

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Infrastructure"
	"bufio"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type ITwitchClient interface {
	ConnectToIrcServer() error
	Disconnect()
	Authenticate(userName string, oauthToken string) error
	JoinChannel(channel string) error
	LeaveChannel(channel string) error
	SendPong() error
	ReadLine() (string, error)
	WriteMessage(message string, channel string, messageType string, user string) (int, error)
}

func NewTwitchClient(settings *Configuration.Settings) ITwitchClient {
	client := new(twitchClient)
	client.Port = settings.TwitchPort
	client.Server = settings.TwitchServer
	client.RateLimiter = Infrastructure.NewRateLimiter(
		settings.TwitchRateLimit, time.Duration(settings.TwitchRateLimitDurationSeconds)*time.Second)
	return client
}

type twitchClient struct {
	Port        string
	Server      string
	RateLimiter Infrastructure.IRateLimiter

	connection net.Conn
	reader     *textproto.Reader
}

func (this *twitchClient) ConnectToIrcServer() error {
	connectionString := this.Server + ":" + this.Port
	log.Println("Connecting to Twitch IRC server at " + connectionString)

	var err error
	this.connection, err = net.Dial("tcp", connectionString)

	if err != nil {
		log.Println("Failed to connect to Twitch IRC server at " + connectionString)
		return err
	}
	log.Println("Successfully connected to Twitch IRC")

	this.reader = textproto.NewReader(bufio.NewReader(this.connection))
	return nil
}

func (this *twitchClient) Authenticate(userName string, oauthToken string) error {

	if _, err := this.writeToConnection("PASS " + oauthToken + "\r\n"); err != nil {
		log.Println("Error passing oauth token to twitch.")
		return err
	}

	if _, err := this.writeToConnection("NICK " + userName + "\r\n"); err != nil {
		log.Println("Error logging in with user: " + userName)
		return err
	}

	if _, err := this.writeToConnection("CAP REQ :twitch.tv/commands\r\n"); err != nil {
		log.Println("Error requesting Twitch command capabilities")
		return err
	}

	log.Println("Successfully authenticated with user: " + userName)
	return nil
}

func (this *twitchClient) JoinChannel(channel string) error {

	if _, err := this.writeToConnection("JOIN #" + strings.ToLower(channel) + "\r\n"); err != nil {
		log.Println("Error joining channel: " + channel)
		return err
	}

	log.Println("Successfully joined channel: " + channel)

	return nil
}

func (this *twitchClient) LeaveChannel(channel string) error {

	if _, err := this.writeToConnection("PART #" + strings.ToLower(channel) + "\r\n"); err != nil {
		log.Println("Error leaving channel: " + channel)
		return err
	}

	log.Println("Successfully left channel: " + channel)

	return nil
}

func (this *twitchClient) Disconnect() {

	err := this.connection.Close()
	if err != nil {
		log.Println("Error disconnecting from Twitch IRC. Error: " + err.Error())
		return
	}

	log.Println("Disconnected from Twitch IRC")
	return
}

func (this *twitchClient) SendPong() error {

	// We do not use this.writeToConnection because we do not want our pong's to be rate limited
	_, err := this.connection.Write([]byte(pongCommand))
	if err != nil {
		log.Println("Error sending PONG!")
		return err
	}

	return nil
}

func (this *twitchClient) ReadLine() (string, error) {
	return this.reader.ReadLine()
}

func (this *twitchClient) WriteMessage(message string, channel string, messageType string, user string) (int, error) {

	// Twitch does not allow for multi-line messages, so we replace any \r\n or \n with a space
	// to make it look nice in chat
	message = strings.ReplaceAll(message, "\r", "")
	message = strings.ReplaceAll(message, "\n", " ")

	command := RegularMessageType + " #" + strings.ToLower(channel) + " :"
	if messageType == WhisperMessageType {
		command += "/w " + user + " "
	}
	command += user + ": " + message + "\r\n"

	var bytesWritten int
	if bytesWritten, err := this.writeToConnection(command); err != nil {
		log.Println("Error writing message to channel: " + channel)
		return bytesWritten, err
	}

	return bytesWritten, nil
}

func (this *twitchClient) writeToConnection(message string) (int, error) {
	var err error = nil
	var bytesWritten int
	this.RateLimiter.PerformInteraction(func() {
		bytesWritten, err = this.connection.Write([]byte(message))
	})
	return bytesWritten, err
}
