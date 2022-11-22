package utils

import (
	"autosign/vo"
	"encoding/json"
)

func GetSubmitForm(userId string, deviceId string) vo.SubmitForm {
	jsonStr := `{
			"lon": "121.640556",
			"version": "first_v3",
			"calVersion": "firstv",
			"deviceId": "",
			"userId": "",
			"systemName": "iOS",
			"bodyString": "",
			"lat": "29.907569",
			"systemVersion": "16.0.2",
			"appVersion": "10.0.1",
			"model": "iPhone 14 Pro Max",
			"sign": "48951d3322c9c507244d2fdbb7280bf4"
	}`

	var form vo.SubmitForm
	err := json.Unmarshal([]byte(jsonStr), &form)
	if err != nil {
		return vo.SubmitForm{}
	}
	form.UserId = userId
	form.DeviceId = deviceId
	return form
}
