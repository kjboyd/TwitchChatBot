package BusinessLogic

type ITwitchBotManager interface {
	RunTwitchBot(done chan bool)
}

func NewTwitchBotManager(twitchChatMonitor ITwitchChatMonitor) ITwitchBotManager {
	manager := new(twitchBotManager)
	manager.TwitchChatMonitor = twitchChatMonitor
	return manager
}

type twitchBotManager struct {
	TwitchChatMonitor ITwitchChatMonitor
}

func (this *twitchBotManager) RunTwitchBot(done chan bool) {
	defer func() {
		this.TwitchChatMonitor.DisconnectFromTwitch()
		done <- true
	}()

	this.TwitchChatMonitor.ConnectToTwitch()

	for running := true; running; {
		running = this.TwitchChatMonitor.MonitorChat()
	}
	return
}
