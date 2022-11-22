package vo

type UploadBody struct {
	Position        string `json:"position"`
	Longitude       string `json:"longitude"`
	Latitude        string `json:"latitude"`
	SignVersion     string `json:"signVersion"`
	UaIsCpadaily    bool   `json:"uaIsCpadaily"`
	SignPhotoUrl    string `json:"signPhotoUrl"`
	IsNeedExtra     int    `json:"isNeedExtra"`
	Ticket          string `json:"ticket"`
	ExtraFieldItems []struct {
		ExtraFieldItemValue string `json:"extraFieldItemValue"`
		ExtraFieldItemWid   int    `json:"extraFieldItemWid"`
	} `json:"extraFieldItems"`
	AbnormalReason  string `json:"abnormalReason"`
	SignInstanceWid string `json:"signInstanceWid"`
	IsMalposition   int    `json:"isMalposition"`
}
