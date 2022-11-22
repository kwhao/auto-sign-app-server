package utils

import (
	"autosign/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

func CheckValidation(cookies []*http.Cookie, deviceId string) (model.ValidateConfig, []*http.Cookie) {
	const validateUrl = "https://nbu.campusphere.net/wec-counselor-sign-apps/stu/sign/checkValidation"
	type body struct {
		QrCode   string `json:"qrCode"`
		DeviceId string `json:"deviceId"`
	}
	validateBody := &body{QrCode: "", DeviceId: deviceId}
	validateForm, _ := json.Marshal(validateBody)
	validateReq, _ := http.NewRequest("POST", validateUrl, strings.NewReader(string(validateForm)))
	validateClient := &http.Client{}
	for _, cookie := range cookies {
		validateReq.AddCookie(cookie)
	}
	validateReq.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4379391488) cpdaily/10.0.1  wisedu/10.0.1")
	validateReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	validateReq.Header.Add("Host", "https://nbu.campusphere.net/")

	res, err := validateClient.Do(validateReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	for _, cookie := range res.Cookies() {
		validateReq.AddCookie(cookie)
	}
	result, _ := ioutil.ReadAll(res.Body)

	var validateConfig model.ValidateConfig
	err = json.Unmarshal(result, &validateConfig)
	fmt.Println("#=============+++++++++++++================#")
	if err != nil {
		return validateConfig, nil
	}
	return validateConfig, validateReq.Cookies()
}

func GetValidateImages(cookies []*http.Cookie, deviceId string, validateInfo model.ValidateInfo) (model.ImageListResult, []*http.Cookie) {
	var imageListResult model.ImageListResult
	const applyValidateUrl = "https://nbu.campusphere.net/captcha-open-api/v1/captcha/getLimitExpire/scenesImage"
	bodyBuf, contentType := getValidateFormBody(validateInfo)
	applyValidateReq, _ := http.NewRequest("POST", applyValidateUrl, bodyBuf)
	applyValidateClient := &http.Client{}
	for _, cookie := range cookies {
		applyValidateReq.AddCookie(cookie)
	}
	applyValidateReq.Header.Set("Content-Type", contentType)
	applyValidateReq.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4379391488) cpdaily/10.0.1  wisedu/10.0.1")
	applyValidateReq.Header.Add("Host", "https://nbu.campusphere.net/")

	applyRes, _ := applyValidateClient.Do(applyValidateReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(applyRes.Body)
	result, _ := ioutil.ReadAll(applyRes.Body)

	var applyResult model.ApplyResult
	_ = json.Unmarshal(result, &applyResult)
	if applyResult.Code != 200 {
		return imageListResult, nil
	}

	const validateImageUrl = "https://nbu.campusphere.net/captcha-open-api/v1/captcha/create/scenesImage"
	bodyBuf, contentType = getValidateFormBody(validateInfo)
	getValidateImageListReq, _ := http.NewRequest("POST", validateImageUrl, bodyBuf)
	getValidateImageListClient := &http.Client{}
	for _, cookie := range cookies {
		getValidateImageListReq.AddCookie(cookie)
	}
	for _, cookie := range applyRes.Cookies() {
		getValidateImageListReq.AddCookie(cookie)
	}
	getValidateImageListReq.Header.Set("Content-Type", contentType)
	getValidateImageListReq.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4379391488) cpdaily/10.0.1  wisedu/10.0.1")
	getValidateImageListReq.Header.Add("Host", "https://nbu.campusphere.net/")
	getValidateImageListReq.Header.Add("deviceId", deviceId)

	imageListRes, _ := getValidateImageListClient.Do(getValidateImageListReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(imageListRes.Body)
	result, _ = ioutil.ReadAll(imageListRes.Body)

	err := json.Unmarshal(result, &imageListResult)
	if err != nil {
		return imageListResult, nil
	}

	// 重制cookie
	for _, cookie := range imageListRes.Cookies() {
		getValidateImageListReq.AddCookie(cookie)
	}
	return imageListResult, getValidateImageListReq.Cookies()

}

func ConfirmCaptcha(cookies []*http.Cookie, deviceId string, validateInfo model.ValidateInfo, captchaBody model.CaptchaBody) (model.CaptchaResult, []*http.Cookie) {
	var captchaResult model.CaptchaResult
	const confirmUrl = "https://nbu.campusphere.net/captcha-open-api/v1/captcha/validate/scenesImage"
	bodyBuf, contentType := getCaptchaFormBody(validateInfo, captchaBody)

	confirmReq, _ := http.NewRequest("POST", confirmUrl, bodyBuf)
	confirmClient := &http.Client{}
	for _, cookie := range cookies {
		confirmReq.AddCookie(cookie)
	}
	confirmReq.Header.Set("Content-Type", contentType)
	confirmReq.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4379391488) cpdaily/10.0.1  wisedu/10.0.1")
	confirmReq.Header.Add("Host", "https://nbu.campusphere.net/")
	confirmReq.Header.Add("deviceId", deviceId)

	confirmRes, _ := confirmClient.Do(confirmReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(confirmRes.Body)
	result, _ := ioutil.ReadAll(confirmRes.Body)
	fmt.Println(string(result))

	err := json.Unmarshal(result, &captchaResult)
	if err != nil {
		return captchaResult, nil
	}

	// 重制cookie
	for _, cookie := range confirmRes.Cookies() {
		confirmReq.AddCookie(cookie)
	}

	return captchaResult, confirmReq.Cookies()

}

func HandleCheatModule(cookies []*http.Cookie) []*http.Cookie {
	const cheatModuleUrl = "https://nbu.campusphere.net/wec-counselor-sign-apps/stu/sign/getCheatModuleStatus"
	cheatReq, _ := http.NewRequest("GET", cheatModuleUrl, nil)
	cheatClient := &http.Client{}
	for _, cookie := range cookies {
		cheatReq.AddCookie(cookie)
	}
	cheatReq.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4379391488) cpdaily/10.0.1  wisedu/10.0.1")
	cheatReq.Header.Add("Host", "https://nbu.campusphere.net/")

	cheatRes, _ := cheatClient.Do(cheatReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(cheatRes.Body)
	result, _ := ioutil.ReadAll(cheatRes.Body)
	fmt.Println(string(result))

	fmt.Println(cheatRes.Header)

	// 重制cookie
	for _, cookie := range cheatRes.Cookies() {
		cheatReq.AddCookie(cookie)
	}

	return cheatReq.Cookies()
}

func getValidateFormBody(validateInfo model.ValidateInfo) (*bytes.Buffer, string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	accountKey, err := bodyWriter.CreateFormField("accountKey")
	if err != nil {
		return nil, ""
	}
	_, errs := accountKey.Write([]byte(validateInfo.AccountKey))
	if errs != nil {
		return nil, ""
	}
	sceneCode, err := bodyWriter.CreateFormField("sceneCode")
	if err != nil {
		return nil, ""
	}
	_, errs = sceneCode.Write([]byte(validateInfo.SceneCode))
	if errs != nil {
		return nil, ""
	}
	tenantId, err := bodyWriter.CreateFormField("tenantId")
	if err != nil {
		return nil, ""
	}
	_, errs = tenantId.Write([]byte(validateInfo.TenantId))
	if errs != nil {
		return nil, ""
	}
	userId, err := bodyWriter.CreateFormField("userId")
	if err != nil {
		return nil, ""
	}
	_, errs = userId.Write([]byte(validateInfo.UserId))
	if errs != nil {
		return nil, ""
	}
	err = bodyWriter.Close()
	if err != nil {
		return nil, ""
	}
	return bodyBuf, bodyWriter.FormDataContentType()
}

func getCaptchaFormBody(validateInfo model.ValidateInfo, captchaBody model.CaptchaBody) (*bytes.Buffer, string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	accountKey, err := bodyWriter.CreateFormField("accountKey")
	if err != nil {
		return nil, ""
	}
	_, errs := accountKey.Write([]byte(validateInfo.AccountKey))
	if errs != nil {
		return nil, ""
	}
	sceneCode, err := bodyWriter.CreateFormField("sceneCode")
	if err != nil {
		return nil, ""
	}
	_, errs = sceneCode.Write([]byte(validateInfo.SceneCode))
	if errs != nil {
		return nil, ""
	}
	tenantId, err := bodyWriter.CreateFormField("tenantId")
	if err != nil {
		return nil, ""
	}
	_, errs = tenantId.Write([]byte(validateInfo.TenantId))
	if errs != nil {
		return nil, ""
	}
	userId, err := bodyWriter.CreateFormField("userId")
	if err != nil {
		return nil, ""
	}
	_, errs = userId.Write([]byte(validateInfo.UserId))
	if errs != nil {
		return nil, ""
	}

	scenesImageCode, err := bodyWriter.CreateFormField("scenesImageCode")
	if err != nil {
		return nil, ""
	}
	_, errs = scenesImageCode.Write([]byte(captchaBody.ScenesImageCode))
	if errs != nil {
		return nil, ""
	}

	for _, singleCode := range captchaBody.ScenesImageCodes {
		scenesImageCodes, err := bodyWriter.CreateFormField("scenesImageCodes")
		if err != nil {
			return nil, ""
		}
		_, errs = scenesImageCodes.Write([]byte(singleCode))
		if errs != nil {
			return nil, ""
		}
	}

	err = bodyWriter.Close()
	if err != nil {
		return nil, ""
	}
	return bodyBuf, bodyWriter.FormDataContentType()
}
