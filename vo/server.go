package vo

type CaptchaBodyVo struct {
	Cookie           string   `json:"cookie"`
	DeviceId         string   `json:"deviceId"`
	AccountKey       string   `json:"accountKey"`
	SceneCode        string   `json:"sceneCode"`
	TenantId         string   `json:"tenantId"`
	UserId           string   `json:"userId"`
	ScenesImageCode  string   `json:"scenesImageCode"`
	ScenesImageCodes []string `json:"scenesImageCodes"`
}
