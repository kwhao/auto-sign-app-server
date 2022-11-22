package vo

type SubmitForm struct {
	Lon           string `json:"lon"`
	Version       string `json:"version"`
	CalVersion    string `json:"calVersion"`
	DeviceId      string `json:"deviceId"`
	UserId        string `json:"userId"`
	SystemName    string `json:"systemName"`
	BodyString    string `json:"bodyString"`
	Lat           string `json:"lat"`
	SystemVersion string `json:"systemVersion"`
	AppVersion    string `json:"appVersion"`
	Model         string `json:"model"`
	Sign          string `json:"sign"`
}
