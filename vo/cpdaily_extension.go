package vo

type CpdailyExtension struct {
	Lon           string `json:"lon"`
	Lat           string `json:"lat"`
	Model         string `json:"model"`
	AppVersion    string `json:"appVersion"`
	SystemVersion string `json:"systemVersion"`
	UserId        string `json:"userId"`
	SystemName    string `json:"systemName"`
	DeviceId      string `json:"deviceId"`
	Ticket        string `json:"ticket"`
}
