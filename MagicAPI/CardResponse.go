package MagicAPI

type CardResponse struct {
	Card  *MagicCard   `json:"card"`
	Cards []*MagicCard `json:"cards"`
}
