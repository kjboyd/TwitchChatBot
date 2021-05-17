package BusinessLogic

func RunChatBot(chatBot IChatBot, done chan bool) {
	defer func() {
		chatBot.Disconnect()
		done <- true
	}()

	chatBot.Connect()

	for running := true; running; {
		running = chatBot.ProcessMessage()
	}
	return
}
