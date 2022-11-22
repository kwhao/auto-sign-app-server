package model

type ExtraFieldItemVo struct {
	FieldIndex            int    `json:"fieldIndex"`
	ExtraTitle            string `json:"extraTitle"`
	ExtraDesc             string `json:"extraDesc"`
	ExtraFieldItemWid     string `json:"extraFieldItemWid"`
	ExtraFieldItem        string `json:"extraFieldItem"`
	IsExtraFieldOtherItem string `json:"isExtraFieldOtherItem"`
	IsAbnormal            string `json:"isAbnormal"`
}
