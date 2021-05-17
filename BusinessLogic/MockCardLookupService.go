package BusinessLogic

type MockCardLookupService struct {
	MethodCalls []string
}

func (this *MockCardLookupService) LookupCardAndPost(
	cardName string, messageType string, channel string, user string) {

	this.MethodCalls = append(this.MethodCalls, "LookupCardAndPost")
}
