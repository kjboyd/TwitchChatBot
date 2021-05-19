package BusinessLogic

import (
	"TwitchChatBot/BusinessLogic/mock_BusinessLogic"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/MagicAPI/mock_MagicAPI"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
)

type cardLookupServiceTestHarness struct {
	MockController *gomock.Controller
	MagicClient    *mock_MagicAPI.MockIMagicClient
	Writer         *mock_BusinessLogic.MockStringWriter
	Patient        ICardLookupService
}

func setupCardLookupServiceTestHarness(test *testing.T) *cardLookupServiceTestHarness {

	testHarness := new(cardLookupServiceTestHarness)
	testHarness.MockController = gomock.NewController(test)
	testHarness.MagicClient = mock_MagicAPI.NewMockIMagicClient(testHarness.MockController)
	testHarness.Writer = mock_BusinessLogic.NewMockStringWriter(testHarness.MockController)
	testHarness.Patient = NewCardLookupService(testHarness.MagicClient)
	return testHarness
}

func Test_WillNotLookupCardIfCardNameIsBlank(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Any()).Times(0)
	testHarness.Writer.EXPECT().WriteString(gomock.Any()).Do(
		func(message string) {
			defer waitGroup.Done()
		},
	)

	testHarness.Patient.LookupCardAndPost("", testHarness.Writer)
	waitGroup.Wait()
}

func Test_WillInformUserWhenCardNameIsBlank(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	testHarness.Writer.EXPECT().WriteString(gomock.Any()).Do(
		func(message string) {
			defer waitGroup.Done()

			if message != "Please specify card name." {
				test.Errorf("Writing incorrect message to chat")
			}
		},
	)

	testHarness.Patient.LookupCardAndPost("", testHarness.Writer)
	waitGroup.Wait()
}

func Test_WillLookupCardInformationWhenCardNameNotBlank(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	cardName := "Bonecrusher Giant"
	expectedCard := anonymousMagicCard()
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Eq(cardName)).Times(1).Return(expectedCard, nil)
	testHarness.Writer.EXPECT().WriteString(gomock.Any()).Do(
		func(actualMessage string) {
			defer waitGroup.Done()

			if actualMessage != expectedCard.String() {
				test.Errorf("Sent wrong message to chat. Got: %s, expected: %s", actualMessage, expectedCard.String())
			}
		},
	)

	testHarness.Patient.LookupCardAndPost(cardName, testHarness.Writer)
	waitGroup.Wait()
}

func Test_WillInformUserIfCardCannotBeFound(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	cardName := "Bonecrusher Giant"
	expectedMessage := "Unable to find card " + cardName
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Eq(cardName)).Times(1).Return(nil, fmt.Errorf("Error"))
	testHarness.Writer.EXPECT().WriteString(gomock.Any()).Do(
		func(actualMessage string) {
			defer waitGroup.Done()

			if actualMessage != expectedMessage {
				test.Errorf("Sent wrong message to chat. Got: %s, expected: %s", actualMessage, expectedMessage)
			}
		},
	)

	testHarness.Patient.LookupCardAndPost(cardName, testHarness.Writer)
	waitGroup.Wait()
}

func anonymousMagicCard() *MagicAPI.MagicCard {
	card := new(MagicAPI.MagicCard)
	card.Name = "Bonecrusher Giant"
	card.ManaCost = "2{R}"
	card.Type = "Creature - Giant"
	card.Text = "Lots of text"
	card.Power = "4"
	card.Toughness = "3"
	return card
}
