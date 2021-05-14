package MagicAPI

import (
	"fmt"
)

// There are more fields then this, but for now, these are the only ones I care about
type MagicCard struct {
	Name         string   `json:"name"`
	Names        []string `json:"names"`
	ManaCost     string   `json:"manaCost"`
	CMC          float64  `json:"cmc"`
	Type         string   `json:"type"`
	Rarity       string   `json:"rarity"`
	SetName      string   `json:"setName"`
	Text         string   `json:"text"`
	Power        string   `json:"power"`
	Toughness    string   `json:"toughness"`
	Loyalty      string   `json:"loyalty"`
	Layout       string   `json:"layout"`
	MultiverseId string   `json:"multiverseid"`
	Id           string   `json:"id"`
}

func (this *MagicCard) String() string {
	displayString := fmt.Sprintf("%s, %s, %s, %s",
		this.Name, this.ManaCost, this.Type, this.Text)

	if this.Power != "" && this.Toughness != "" {
		displayString += fmt.Sprintf(", %s/%s", this.Power, this.Toughness)
	}

	return displayString
}
