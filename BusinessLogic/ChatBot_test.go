package BusinessLogic

import (
	"TwitchChatBot/Configuration"
	"TwitchChatBot/TwitchAPI"
	"testing"
)

type chatBotTestHarness struct {
	CardLookupService *MockCardLookupService
	TwitchClient      *TwitchAPI.MockTwitchClient
	Settings          *Configuration.Settings
	Patient           IChatBot
}

func setupTestHarness() *chatBotTestHarness {
	testHarness := new(chatBotTestHarness)
	testHarness.CardLookupService = new(MockCardLookupService)
	testHarness.TwitchClient = new(TwitchAPI.MockTwitchClient)
	testHarness.Settings = new(Configuration.Settings)
	testHarness.Patient = NewChatBot(testHarness.TwitchClient, testHarness.CardLookupService, testHarness.Settings)
	return testHarness
}

func Test_WillConnectToTwitchCorrectly(test *testing.T) {
	testHarness := setupTestHarness()

	testHarness.Patient.Connect()

	// In the real world I will use a mocking framework so that I
	// can test arguments and do these checks as expectations rather
	// than checking items in a list
	if len(testHarness.TwitchClient.MethodCalls) != 4 {
		test.Errorf("Connect to Irc Server was never called.")
	}

	if testHarness.TwitchClient.MethodCalls[0] != "ConnectToIrcServer" {
		test.Errorf("Connect to Irc Server was never called.")
	}

	if testHarness.TwitchClient.MethodCalls[1] != "Authenticate" {
		test.Errorf("Authenticate was never called.")
	}

	if testHarness.TwitchClient.MethodCalls[2] != "JoinChannel" {
		test.Errorf("JoinChannel was never called.")
	}

	if testHarness.TwitchClient.MethodCalls[3] != "JoinChannel" {
		test.Errorf("JoinChannel was never called.")
	}
}
