package BusinessLogic

import "TwitchChatBot/TwitchAPI"

func NewTwitchWriter(twitchClient TwitchAPI.ITwitchClient, messageType string, channel string, user string) *twitchWriter {
	writer := new(twitchWriter)
	writer.twitchClient = twitchClient
	writer.messageType = messageType
	writer.channel = channel
	writer.user = user
	return writer
}

// Implements interface StringWriter
type twitchWriter struct {
	twitchClient TwitchAPI.ITwitchClient
	messageType  string
	channel      string
	user         string
}

func (this *twitchWriter) WriteString(str string) (int, error) {
	return this.twitchClient.WriteMessage(str, this.channel, this.messageType, this.user)
}
