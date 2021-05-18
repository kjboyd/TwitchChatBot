package BusinessLogic

import (
	"TwitchChatBot/BusinessLogic/mock_BusinessLogic"
	"TwitchChatBot/Configuration"
	"TwitchChatBot/TwitchAPI/mock_TwitchAPI"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

type chatBotTestHarness struct {
	CardLookupService *mock_BusinessLogic.MockICardLookupService
	TwitchClient      *mock_TwitchAPI.MockITwitchClient
	Settings          *Configuration.Settings
	Patient           IChatBot
}

func setupTestHarness(test *testing.T) *chatBotTestHarness {
	ctrl := gomock.NewController(test)

	testHarness := new(chatBotTestHarness)
	testHarness.CardLookupService = mock_BusinessLogic.NewMockICardLookupService(ctrl)
	testHarness.TwitchClient = mock_TwitchAPI.NewMockITwitchClient(ctrl)
	testHarness.Settings = new(Configuration.Settings)
	testHarness.Settings.UserName = "testUser"
	testHarness.Settings.AuthToken = "testAuth"
	testHarness.Settings.Channel = "testChannel"
	testHarness.Patient = NewChatBot(testHarness.TwitchClient, testHarness.CardLookupService, testHarness.Settings)
	return testHarness
}

func Test_WillConnectToTwitchCorrectly(test *testing.T) {
	testHarness := setupTestHarness(test)

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer()
	testHarness.TwitchClient.EXPECT().Authenticate(
		gomock.Eq(testHarness.Settings.UserName),
		gomock.Eq(testHarness.Settings.AuthToken))
	testHarness.TwitchClient.EXPECT().JoinChannel(
		gomock.Eq(testHarness.Settings.UserName))
	testHarness.TwitchClient.EXPECT().JoinChannel(
		gomock.Eq(testHarness.Settings.Channel))

	testHarness.Patient.Connect()
}

func Test_WillReturnErrorIfConnectToIrcServerFails(test *testing.T) {
	testHarness := setupTestHarness(test)

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer().Return(
		fmt.Errorf("Error"))

	err := testHarness.Patient.Connect()

	if err == nil {
		test.Errorf("Connect succeeded even though connect to IrcServer Failed")
	}
}
