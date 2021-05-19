package BusinessLogic

import (
	"TwitchChatBot/BusinessLogic/mock_BusinessLogic"
	"TwitchChatBot/MagicAPI"
	"TwitchChatBot/MagicAPI/mock_MagicAPI"
	"fmt"
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

	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Any()).Times(0)
	testHarness.Writer.EXPECT().WriteString(gomock.Any())

	testHarness.Writer.WaitGroup.Add(1)
	testHarness.Patient.LookupCardAndPost("", testHarness.Writer)
	testHarness.Writer.WaitGroup.Wait()
}

func Test_WillInformUserWhenCardNameIsBlank(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	testHarness.Writer.EXPECT().WriteString("Please specify card name.")

	testHarness.Writer.WaitGroup.Add(1)
	testHarness.Patient.LookupCardAndPost("", testHarness.Writer)
	testHarness.Writer.WaitGroup.Wait()
}

func Test_WillLookupCardInformationWhenCardNameNotBlank(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	cardName := "Bonecrusher Giant"
	expectedCard := anonymousMagicCard()
	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Eq(cardName)).Times(1).Return(expectedCard, nil)
	testHarness.Writer.EXPECT().WriteString(expectedCard.String())

	testHarness.Writer.WaitGroup.Add(1)
	testHarness.Patient.LookupCardAndPost(cardName, testHarness.Writer)
	testHarness.Writer.WaitGroup.Wait()
}

func Test_WillInformUserIfCardCannotBeFound(test *testing.T) {
	testHarness := setupCardLookupServiceTestHarness(test)
	defer testHarness.MockController.Finish()

	cardName := "Bonecrusher Giant"
	testHarness.MagicClient.EXPECT().LookupCardInformation(
		gomock.Eq(cardName)).Times(1).Return(nil, fmt.Errorf("Error"))
	testHarness.Writer.EXPECT().WriteString("Unable to find card " + cardName)

	testHarness.Writer.WaitGroup.Add(1)
	testHarness.Patient.LookupCardAndPost(cardName, testHarness.Writer)
	testHarness.Writer.WaitGroup.Wait()
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
