package model

type ExtraFieldItem struct {
	Content      string `json:"content"`
	Wid          int    `json:"wid"`
	IsOtherItems int    `json:"isOtherItems"`
	Value        string `json:"value"`
	IsSelected   bool   `json:"isSelected"`
	IsAbnormal   bool   `json:"isAbnormal"`
}
