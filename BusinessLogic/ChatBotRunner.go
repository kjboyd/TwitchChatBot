package BusinessLogic

func RunChatBot(chatBot IChatBot, done chan bool) {
	defer func() {
		chatBot.Disconnect()
		done <- true
	}()

	err := chatBot.Connect()
	if err != nil {
		return
	}

	for running := true; running; {
		running = chatBot.ProcessMessage()
	}
	return
}
