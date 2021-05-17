package TwitchAPI

type MockTwitchClient struct {
	ConnectToIrcServerError error
	AuthenticateError       error
	JoinChannelError        error
	LeaveChannelError       error
	SendPongError           error
	ReadLineResponse        string
	ReadLineError           error
	WriteMessageError       error
	MethodCalls             []string
}

func (this *MockTwitchClient) ConnectToIrcServer() error {
	this.MethodCalls = append(this.MethodCalls, "ConnectToIrcServer")
	return this.ConnectToIrcServerError
}

func (this *MockTwitchClient) Disconnect() {
	this.MethodCalls = append(this.MethodCalls, "Disconnect")
}

func (this *MockTwitchClient) Authenticate(userName string, oauthToken string) error {
	this.MethodCalls = append(this.MethodCalls, "Authenticate")
	return this.AuthenticateError
}

func (this *MockTwitchClient) JoinChannel(channel string) error {
	this.MethodCalls = append(this.MethodCalls, "JoinChannel")
	return this.JoinChannelError
}

func (this *MockTwitchClient) LeaveChannel(channel string) error {
	this.MethodCalls = append(this.MethodCalls, "LeaveChannel")
	return this.LeaveChannelError
}

func (this *MockTwitchClient) SendPong() error {
	this.MethodCalls = append(this.MethodCalls, "SendPong")
	return this.SendPongError
}

func (this *MockTwitchClient) ReadLine() (string, error) {
	this.MethodCalls = append(this.MethodCalls, "ReadLine")
	return this.ReadLineResponse, this.ReadLineError
}

func (this *MockTwitchClient) WriteMessage(message string, channel string, messageType string, user string) error {
	this.MethodCalls = append(this.MethodCalls, "WriteMessage")
	return this.WriteMessageError
}
