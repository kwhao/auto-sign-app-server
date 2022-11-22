package model

type ExtraField struct {
	Wid             int              `json:"wid"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	HasOtherItems   int              `json:"hasOtherItems"`
	ExtraFieldItems []ExtraFieldItem `json:"extraFieldItems"`
}
