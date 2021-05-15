package TwitchAPI

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/Infrastructure"
	"TwitchChatBot/Logging"
	"bufio"
	"net"
	"net/textproto"
	"strings"
	"time"
)

const pong string = "PONG :tmi.twitch.tv\r\n"

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

func NewTwitchClient(settings *Configuration.Settings, logger Logging.ILogger) ITwitchClient {
	client := new(twitchClient)
	client.Port = settings.TwitchPort
	client.Server = settings.TwitchServer
	client.Logger = logger
	client.RateLimiter = Infrastructure.NewRateLimiter(
		settings.TwitchRateLimit, time.Duration(settings.TwitchRateLimitDurationSeconds)*time.Second)
	return client
}

type twitchClient struct {
	Port        string
	Server      string
	Logger      Logging.ILogger
	RateLimiter Infrastructure.IRateLimiter

	connection net.Conn
	reader     *textproto.Reader
}

func (this *twitchClient) ConnectToIrcServer() error {
	connectionString := this.Server + ":" + this.Port
	this.Logger.Log("Connecting to Twitch IRC server at " + connectionString)

	var err error
	this.connection, err = net.Dial("tcp", connectionString)

	if err != nil {
		this.Logger.Log("Failed to connect to Twitch IRC server at " + connectionString)
		return err
	}
	this.Logger.Log("Successfully connected to Twitch IRC")

	this.reader = textproto.NewReader(bufio.NewReader(this.connection))
	return nil
}

func (this *twitchClient) Authenticate(userName string, oauthToken string) error {

	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("PASS " + oauthToken + "\r\n"))
	})
	if err != nil {
		this.Logger.Log("Error passing oauth token to twitch.")
		return err
	}

	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("NICK " + userName + "\r\n"))
	})
	if err != nil {
		this.Logger.Log("Error logging in with user: " + userName)
		return err
	}

	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("CAP REQ :twitch.tv/commands\r\n"))
	})
	if err != nil {
		this.Logger.Log("Error requesting commands")
		return err
	}

	this.Logger.Log("Successfully authenticated with user: " + userName)
	return nil
}

func (this *twitchClient) JoinChannel(channel string) error {

	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("JOIN #" + strings.ToLower(channel) + "\r\n"))
	})

	if err != nil {
		this.Logger.Log("Error joining channel: " + channel)
		return err
	}

	this.Logger.Log("Successfully joined channel: " + channel)

	return nil
}

func (this *twitchClient) LeaveChannel(channel string) error {
	var err error = nil
	this.RateLimiter.PerformInteraction(func() {
		_, err = this.connection.Write([]byte("PART #" + strings.ToLower(channel) + "\r\n"))
	})

	if err != nil {
		this.Logger.Log("Error leaving channel: " + channel)
		return err
	}

	this.Logger.Log("Successfully left channel: " + channel)

	return nil
}

func (this *twitchClient) Disconnect() {

	err := this.connection.Close()
	if err != nil {
		this.Logger.Log("Error disconnecting from Twitch IRC")
		return
	}

	this.Logger.Log("Disconnected from Twitch IRC")
	return
}

func (this *twitchClient) SendPong() error {

	_, err := this.connection.Write([]byte(pong))
	if err != nil {
		this.Logger.Log("Error sending PONG!")
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

		command := "PRIVMSG #" + strings.ToLower(channel) + " :"
		if messageType == "WHISPER" {
			command += "/w " + user + " "
		}
		command += user + ": " + message + "\r\n"

		_, err = this.connection.Write([]byte(command))
		if err != nil {
			this.Logger.Log("Error writing message to channel: " + channel)
		}
	})

	return err
}
