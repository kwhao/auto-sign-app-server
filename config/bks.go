package config

import (
	"autosign/vo"
	"encoding/json"
)

func GetBksConfig(validateTicket string) vo.UploadBody {
	jsonStr := `{
    "position": "浙江省宁波市江北区宁波大学",
    "longitude": "121.640556",
    "latitude": "29.907569",
    "signVersion": "1.0.0",
    "uaIsCpadaily": true,
    "signPhotoUrl": "",
    "isNeedExtra": 1,
	"ticket":"",
    "extraFieldItems": [
        {
            "extraFieldItemValue": "健康"
        },
        {
            "extraFieldItemValue": "绿色"
        },
        {
            "extraFieldItemValue": "无"
        },
        {
            "extraFieldItemValue": "是，在校住宿"
        }
    ],
    "abnormalReason": "实习",
    "isMalposition": 0
}`
	var form vo.UploadBody
	err := json.Unmarshal([]byte(jsonStr), &form)
	if err != nil {
		return vo.UploadBody{}
	}
	form.Ticket = validateTicket
	return form
}
