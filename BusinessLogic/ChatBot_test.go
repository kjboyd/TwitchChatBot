package BusinessLogic

import (
	"TwitchChatBot/BusinessLogic/mock_BusinessLogic"
	"TwitchChatBot/Configuration"
	"TwitchChatBot/TwitchAPI"
	"TwitchChatBot/TwitchAPI/mock_TwitchAPI"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
)

type chatBotTestHarness struct {
	MockController    *gomock.Controller
	CardLookupService *mock_BusinessLogic.MockICardLookupService
	TwitchClient      *mock_TwitchAPI.MockITwitchClient
	Settings          *Configuration.Settings
	Patient           IChatBot
}

func setupChatBotTestHarness(test *testing.T) *chatBotTestHarness {

	testHarness := new(chatBotTestHarness)
	testHarness.MockController = gomock.NewController(test)
	testHarness.CardLookupService = mock_BusinessLogic.NewMockICardLookupService(testHarness.MockController)
	testHarness.TwitchClient = mock_TwitchAPI.NewMockITwitchClient(testHarness.MockController)
	testHarness.Settings = new(Configuration.Settings)
	testHarness.Settings.UserName = "testUser"
	testHarness.Settings.AuthToken = "testAuth"
	testHarness.Settings.Channel = "testChannel"
	testHarness.Settings.CardCommand = "card"
	testHarness.Settings.ChangeChannelCommand = "changechannel"
	testHarness.Settings.DisconnectCommand = "shutdown"
	testHarness.Patient = NewChatBot(testHarness.TwitchClient, testHarness.CardLookupService, testHarness.Settings)
	return testHarness
}

func Test_WillConnectToTwitchCorrectly(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

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
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer().Return(
		fmt.Errorf("Error"))

	err := testHarness.Patient.Connect()

	if err == nil {
		test.Errorf("Connect succeeded even though connect to IrcServer failed")
	}
}

func Test_WillReturnErrorIfAuthenticateFails(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer()
	testHarness.TwitchClient.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(
		fmt.Errorf("Error"))

	err := testHarness.Patient.Connect()

	if err == nil {
		test.Errorf("Connect succeeded even though authenticate failed")
	}
}

func Test_WillReturnErrorIfJoinUserChannelFails(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer()
	testHarness.TwitchClient.EXPECT().Authenticate(gomock.Any(), gomock.Any())
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Eq(testHarness.Settings.UserName)).Return(
		fmt.Errorf("Error"))

	err := testHarness.Patient.Connect()

	if err == nil {
		test.Errorf("Connect succeeded even though join user channel failed")
	}
}

func Test_WillReturnErrorIfJoinMainChannelFails(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ConnectToIrcServer()
	testHarness.TwitchClient.EXPECT().Authenticate(gomock.Any(), gomock.Any())
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Eq(testHarness.Settings.UserName))
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Eq(testHarness.Settings.Channel)).Return(
		fmt.Errorf("Error"))

	err := testHarness.Patient.Connect()

	if err == nil {
		test.Errorf("Connect succeeded even though join main channel failed")
	}
}

func Test_WillLeaveAllChannelsAndDisconnectFromTwitchWhenDisconnecting(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().LeaveChannel(
		gomock.Eq(testHarness.Settings.Channel))
	testHarness.TwitchClient.EXPECT().LeaveChannel(
		gomock.Eq(testHarness.Settings.UserName))
	testHarness.TwitchClient.EXPECT().Disconnect()

	testHarness.Patient.Disconnect()
}

func Test_WillStopProcessingMessagesIfUnableToReadLines(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ReadLine().Return("", fmt.Errorf("Error"))

	result := testHarness.Patient.ProcessMessage()

	if result {
		test.Errorf("Continuing to process messages when unable to read lines.")
	}
}

func Test_WillSendPongInResponseToSentPing(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	testHarness.TwitchClient.EXPECT().ReadLine().Return(
		TwitchAPI.Ping, nil)
	testHarness.TwitchClient.EXPECT().SendPong().Do(
		func() {
			waitGroup.Done()
		},
	)

	result := testHarness.Patient.ProcessMessage()
	waitGroup.Wait()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillDoNothingInResponseToUnknownChatMessage(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.TwitchClient.EXPECT().ReadLine().Return(
		"Something strange", nil)

	result := testHarness.Patient.ProcessMessage()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillLookupCardAndPostInResponseToCardCommand(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.RegularMessageType
	channel := "testChannel"
	cardName := "Bonecrusher Giant"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.CardCommand + " " + cardName
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	var writer *twitchWriter
	testHarness.CardLookupService.EXPECT().LookupCardAndPost(gomock.Any(),
		gomock.AssignableToTypeOf(writer)).DoAndReturn(
		func(actualCardName string, writer io.StringWriter) {
			defer waitGroup.Done()

			twitchWriter, _ := writer.(*twitchWriter)
			if actualCardName != cardName {
				test.Errorf("Looking up wrong card name. Got: %s, expected: %s", actualCardName, cardName)
			}
			if twitchWriter.channel != channel {
				test.Errorf("Created twitch writer with wrong channel. Got: %s, expected: %s", twitchWriter.channel, channel)
			}
			if twitchWriter.user != user {
				test.Errorf("Created twitch writer with wrong user. Got: %s, expected: %s", twitchWriter.user, user)
			}
			if twitchWriter.messageType != messageType {
				test.Errorf("Created twitch writer with wrong message type. Got: %s, expected: %s", twitchWriter.messageType, messageType)
			}
			return
		},
	)

	result := testHarness.Patient.ProcessMessage()
	waitGroup.Wait()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillNotStopProcessingMessagesWhenShutdownCommandReceivedAsRegularChatMessage(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.RegularMessageType
	channel := "testChannel"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.DisconnectCommand
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	result := testHarness.Patient.ProcessMessage()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillStopProcessingMessagesWhenShutdownCommandReceivedAsWhisper(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.WhisperMessageType
	channel := "testChannel"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.DisconnectCommand
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	result := testHarness.Patient.ProcessMessage()

	if result {
		test.Errorf("Still processing messages after receiving disconnect command")
	}
}

func Test_WillNotChangeChannelWhenCommandReceivedAsRegularChatMessage(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.RegularMessageType
	channel := "testChannel"
	newChannel := "testChannel2"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.ChangeChannelCommand + " " + newChannel
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	testHarness.TwitchClient.EXPECT().LeaveChannel(gomock.Any()).Times(0)
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Any()).Times(0)

	result := testHarness.Patient.ProcessMessage()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillChangeChannelWhenCommandReceivedAsWhisper(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.WhisperMessageType
	channel := "testChannel"
	newChannel := "testChannel2"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.ChangeChannelCommand + " " + newChannel
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	testHarness.TwitchClient.EXPECT().LeaveChannel(gomock.Any()).Do(
		func(actualChannel string) {
			defer waitGroup.Done()

			if actualChannel != testHarness.Settings.Channel {
				test.Errorf("Leaving wrong channel. Got: %s, expected: %s", actualChannel, testHarness.Settings.Channel)
			}
		},
	)
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Any()).Do(
		func(actualChannel string) {
			defer waitGroup.Done()

			if actualChannel != newChannel {
				test.Errorf("Joining wrong channel. Got: %s, expected: %s", actualChannel, newChannel)
			}
		},
	)

	result := testHarness.Patient.ProcessMessage()
	waitGroup.Wait()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillNotJoinChannelWhenSwitchingToOurOwnChannel(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.WhisperMessageType
	channel := "testChannel"
	newChannel := testHarness.Settings.UserName
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.ChangeChannelCommand + " " + newChannel
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	testHarness.TwitchClient.EXPECT().LeaveChannel(gomock.Any()).Do(
		func(actualChannel string) {
			defer waitGroup.Done()

			if actualChannel != testHarness.Settings.Channel {
				test.Errorf("Leaving wrong channel. Got: %s, expected: %s", actualChannel, testHarness.Settings.Channel)
			}
		},
	)
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Any()).Times(0)

	result := testHarness.Patient.ProcessMessage()
	waitGroup.Wait()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}

func Test_WillNotLeaveChannelWhenSwitchingAwayFromOurOwnChannel(test *testing.T) {
	testHarness := setupChatBotTestHarness(test)
	defer testHarness.MockController.Finish()

	user := "testUser"
	messageType := TwitchAPI.WhisperMessageType
	channel := testHarness.Settings.UserName
	newChannel := "testChannel2"
	chatMessage := ":" + user + "!angelzwrath11@angelzwrath11.tmi.twitch.tv " + messageType + " #" + channel + " :!" + testHarness.Settings.ChangeChannelCommand + " " + newChannel
	testHarness.TwitchClient.EXPECT().ReadLine().Return(chatMessage, nil)
	concretePatient, _ := testHarness.Patient.(*chatBot)
	concretePatient.CurrentChannel = testHarness.Settings.UserName

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	testHarness.TwitchClient.EXPECT().LeaveChannel(gomock.Any()).Times(0)
	testHarness.TwitchClient.EXPECT().JoinChannel(gomock.Any()).Do(
		func(actualChannel string) {
			defer waitGroup.Done()

			if actualChannel != newChannel {
				test.Errorf("Joining wrong channel. Got: %s, expected: %s", actualChannel, newChannel)
			}
		},
	)

	result := testHarness.Patient.ProcessMessage()
	waitGroup.Wait()

	if !result {
		test.Errorf("Stopping processing messages")
	}
}
