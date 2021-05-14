package TwitchAPI

import (
	"TwitchChatBot/Logging"
	"bufio"
	"net"
	"net/textproto"
	"strings"
)

const pong string = "PONG :tmi.twitch.tv\r\n"

type ITwitchClient interface {
	ConnectToIrcServer() error
	Disconnect()
	Authenticate(userName string, oauthToken string) error
	JoinChannel(channel string) error
	SendPong() error
	ReadLine() (string, error)
	WriteMessage(message string, channel string, messageType string, user string) error
}

type TwitchClient struct {
	Port   string
	Server string
	Logger Logging.ILogger

	connection net.Conn
	reader     *textproto.Reader
}

func (this *TwitchClient) ConnectToIrcServer() error {
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

func (this *TwitchClient) Authenticate(userName string, oauthToken string) error {
	_, err := this.connection.Write([]byte("PASS " + oauthToken + "\r\n"))
	if err != nil {
		this.Logger.Log("Error passing oauth token to twitch.")
		return err
	}

	_, err = this.connection.Write([]byte("NICK " + userName + "\r\n"))
	if err != nil {
		this.Logger.Log("Error logging in with user: " + userName)
		return err
	}

	_, err = this.connection.Write([]byte("CAP REQ :twitch.tv/commands\r\n"))
	if err != nil {
		this.Logger.Log("Error requesting commands")
		return err
	}

	this.Logger.Log("Successfully authenticated with user: " + userName)
	return nil
}

func (this *TwitchClient) JoinChannel(channel string) error {
	_, err := this.connection.Write([]byte("JOIN #" + strings.ToLower(channel) + "\r\n"))
	if err != nil {
		this.Logger.Log("Error joining channel: " + channel)
	}

	this.Logger.Log("Successfully joined channel: " + channel)

	return nil
}

func (this *TwitchClient) Disconnect() {
	err := this.connection.Close()
	if err != nil {
		this.Logger.Log("Error disconnecting from Twitch IRC")
		return
	}

	this.Logger.Log("Disconnected from Twitch IRC")
	return
}

func (this *TwitchClient) SendPong() error {
	_, err := this.connection.Write([]byte(pong))
	if err != nil {
		this.Logger.Log("Error sending PONG!")
		return err
	}

	return nil
}

func (this *TwitchClient) ReadLine() (string, error) {
	return this.reader.ReadLine()
}

func (this *TwitchClient) WriteMessage(message string, channel string, messageType string, user string) error {
	command := "PRIVMSG #" + strings.ToLower(channel) + " :"
	if messageType == "WHISPER" {
		command += "/w " + user + " "
	}

	command += message + "\r\n"

	_, err := this.connection.Write([]byte(command))

	if err != nil {
		this.Logger.Log("Error writing message to channel: " + channel)
		return err
	}

	return nil
}
