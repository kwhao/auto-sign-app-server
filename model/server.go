package model

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ValidateConfig struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Validation bool   `json:"validation"`
		AccountKey string `json:"accountKey"`
		SceneCode  string `json:"sceneCode"`
		TenantId   string `json:"tenantId"`
		UserId     string `json:"userId"`
	} `json:"datas"`
}

type ValidateInfo struct {
	AccountKey string `json:"accountKey"`
	SceneCode  string `json:"sceneCode"`
	TenantId   string `json:"tenantId"`
	UserId     string `json:"userId"`
}

type ApplyResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  int    `json:"result"`
}

type ImageListResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Code       string `json:"code"`
		ImageInfos []struct {
			Code string `json:"code"`
			Path string `json:"path"`
		} `json:"imageInfos"`
		Name string `json:"name"`
	} `json:"result"`
}

type CaptchaBody struct {
	ScenesImageCode  string   `json:"scenesImageCode"`
	ScenesImageCodes []string `json:"scenesImageCodes"`
}

type CaptchaResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type AutoSignResult struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Datas   interface{} `json:"datas"`
}
