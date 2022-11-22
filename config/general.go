package config

import (
	"autosign/vo"
	"encoding/json"
)

func GetExtensionInfo(userId string, deviceId string, validateTicket string) vo.CpdailyExtension {
	jsonStr := `
	{"lon": "121.640556",
    "lat": "29.907569",
    "model": "iPhone 14 Pro Max",
    "appVersion": "10.0.1",
    "systemVersion": "16.0.2",
    "userId": "",
    "systemName": "iOS",
	"ticket":"",
    "deviceId": "A2A59D77-3249-4D87-8494-86EFAB53A636"}`

	var form vo.CpdailyExtension
	err := json.Unmarshal([]byte(jsonStr), &form)
	if err != nil {
		return vo.CpdailyExtension{}
	}
	form.UserId = userId
	form.DeviceId = deviceId
	form.Ticket = validateTicket
	//form.DeviceId = strings.ToUpper(uuid.NewV4().String())
	return form
}
