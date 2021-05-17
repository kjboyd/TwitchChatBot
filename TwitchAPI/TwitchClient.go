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
	WriteMessage(message string, channel string, messageType string, user string) error
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

	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("PASS " + oauthToken + "\r\n"))
	})
	if err != nil {
		log.Println("Error passing oauth token to twitch.")
		return err
	}

	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("NICK " + userName + "\r\n"))
	})
	if err != nil {
		log.Println("Error logging in with user: " + userName)
		return err
	}

	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("CAP REQ :twitch.tv/commands\r\n"))
	})
	if err != nil {
		log.Println("Error requesting commands")
		return err
	}

	log.Println("Successfully authenticated with user: " + userName)
	return nil
}

func (this *twitchClient) JoinChannel(channel string) error {

	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("JOIN #" + strings.ToLower(channel) + "\r\n"))
	})

	if err != nil {
		log.Println("Error joining channel: " + channel)
		return err
	}

	log.Println("Successfully joined channel: " + channel)

	return nil
}

func (this *twitchClient) LeaveChannel(channel string) error {
	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("PART #" + strings.ToLower(channel) + "\r\n"))
	})

	if err != nil {
		log.Println("Error leaving channel: " + channel)
		return err
	}

	log.Println("Successfully left channel: " + channel)

	return nil
}

func (this *twitchClient) Disconnect() {

	err := this.connection.Close()
	if err != nil {
		log.Println("Error disconnecting from Twitch IRC")
		return
	}

	log.Println("Disconnected from Twitch IRC")
	return
}

func (this *twitchClient) SendPong() error {

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

func (this *twitchClient) WriteMessage(message string, channel string, messageType string, user string) error {

	var err error = nil
	this.RateLimiter.PerformInteraction(func() {

		// Twitch does not allow for multi-line messages, so we replace any \r\n or \n with a space
		// to make it look nice in chat
		message = strings.ReplaceAll(message, "\r", "")
		message = strings.ReplaceAll(message, "\n", " ")

		command := RegularMessageType + " #" + strings.ToLower(channel) + " :"
		if messageType == WhisperMessageType {
			command += "/w " + user + " "
		}
		command += user + ": " + message + "\r\n"

		_, err = this.connection.Write([]byte(command))
		if err != nil {
			log.Println("Error writing message to channel: " + channel)
		}
	})

	return err
}
