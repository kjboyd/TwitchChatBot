package MagicAPI

type MockMagicClient struct {
	LookupCardInformationError    error
	LookupCardInformationResponse *MagicCard
	MethodCalls                   []string
}

func (this *MockMagicClient) LookupCardInformation(cardNameOrId string) (*MagicCard, error) {
	this.MethodCalls = append(this.MethodCalls, "LookupCardInformation")
	return this.LookupCardInformationResponse, this.LookupCardInformationError
}
