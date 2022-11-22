package main

import (
	"autosign/config"
	"autosign/model"
	"autosign/utils"
	"autosign/vo"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func main() {
	//fmt.Println(config.GetBksConfig())
	fmt.Println("#=============+++++++++++++================#")
	fmt.Println("欢迎使用宁波大学异地在校健康打卡小工具")
	fmt.Println("本工具由Go语言开发")

	r := gin.Default()

	r.GET("/uuid", func(c *gin.Context) {
		c.JSON(200, gin.H{"uuid": strings.ToUpper(uuid.NewV4().String())})
	})

	// 校验账号是否可以验证码登录
	r.GET("/validate", func(c *gin.Context) {
		username := c.Query("username")
		if len(username) == 0 {
			c.JSON(403, gin.H{"msg": "参数不全"})
			return
		}
		// 验证账号是否需要验证码
		if !utils.ValidateAccount(username) {
			fmt.Println("您的账号被限制需要验证码，请先前往网上办事大厅登录一遍后才能使用")
			c.JSON(403, gin.H{"msg": "你的账号被限制需要验证码，请先前往网上办事大厅登录一遍后才能使用"})
			return
		}
		c.JSON(200, gin.H{"msg": "您的账号可以直接登录"})
	})

	r.GET("/login", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if len(username) == 0 || len(password) == 0 {
			c.JSON(403, gin.H{"msg": "参数不全"})
			return
		}

		// 获取登录首页url
		var OpJar *cookiejar.Jar
		OpJar, _ = cookiejar.New(nil)
		client := resty.NewWithClient(
			&http.Client{
				Jar: OpJar,
			},
		)
		resp, err := client.R().
			SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Apple Mac OS X 13_1_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.88 Safari/605.1.15").
			Get("https://nbu.campusphere.net/portal/login")
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error(), "cookie": nil})
			return
		}
		// 进行登录，并获取cookie
		cookies, err := utils.Login(username, password, resp.RawResponse.Request.URL.String())
		cookieList := make([]model.Cookie, 0)
		for _, cookie := range cookies {
			cookieList = append(cookieList, model.Cookie{Name: cookie.Name, Value: cookie.Value})
		}
		cookieJson, _ := json.Marshal(cookieList)
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error(), "cookie": nil})
			return
		}
		c.JSON(200, gin.H{"msg": "success", "cookie": cookieJson})
	})

	r.GET("/checkValidate", func(c *gin.Context) {
		cookieRaw := c.Query("cookie")
		deviceId := c.Query("deviceId")

		if len(cookieRaw) == 0 || len(deviceId) == 0 {
			c.JSON(403, gin.H{"msg": "参数不全"})
			return
		}

		cookieStr, _ := base64.StdEncoding.DecodeString(cookieRaw)

		var cookies []*http.Cookie
		err := json.Unmarshal(cookieStr, &cookies)
		if err != nil {
			c.JSON(405, gin.H{"msg": "cookie无法解析或被篡改"})
			return
		}

		validateConfig, newCookies := utils.CheckValidation(cookies, deviceId)
		cookieList := make([]model.Cookie, 0)
		for _, cookie := range newCookies {
			cookieList = append(cookieList, model.Cookie{Name: cookie.Name, Value: cookie.Value})
		}
		cookieJson, _ := json.Marshal(cookieList)

		if validateConfig.Datas.Validation {
			c.JSON(200, gin.H{"msg": "success", "isNeed": true, "data": validateConfig.Datas, "cookie": cookieJson})
			return
		} else if validateConfig.Message == "SUCCESS" {
			c.JSON(200, gin.H{"msg": "success", "isNeed": false, "data": nil, "cookie": cookieJson})
			return
		} else {
			c.JSON(403, gin.H{"msg": "发生了未知错误", "isNeed": true, "data": nil})
			return
		}
	})

	r.GET("/getValidateImages", func(c *gin.Context) {
		cookieRaw := c.Query("cookie")
		deviceId := c.Query("deviceId")
		accountKey := c.Query("accountKey")
		sceneCode := c.Query("sceneCode")
		tenantId := c.Query("tenantId")
		userId := c.Query("userId")
		if len(cookieRaw) == 0 || len(deviceId) == 0 || len(accountKey) == 0 || len(sceneCode) == 0 || len(tenantId) == 0 || len(userId) == 0 {
			c.JSON(403, gin.H{"msg": "参数不全"})
			return
		}

		cookieStr, _ := base64.StdEncoding.DecodeString(cookieRaw)

		var cookies []*http.Cookie
		err := json.Unmarshal(cookieStr, &cookies)
		if err != nil {
			c.JSON(405, gin.H{"msg": "cookie无法解析或被篡改"})
			return
		}

		validateInfo := model.ValidateInfo{
			AccountKey: accountKey,
			SceneCode:  sceneCode,
			TenantId:   tenantId,
			UserId:     userId,
		}
		imageListResult, newCookies := utils.GetValidateImages(cookies, deviceId, validateInfo)
		cookieList := make([]model.Cookie, 0)
		for _, cookie := range newCookies {
			cookieList = append(cookieList, model.Cookie{Name: cookie.Name, Value: cookie.Value})
		}
		cookieJson, _ := json.Marshal(cookieList)

		if imageListResult.Code == 200 {
			c.JSON(200, gin.H{"msg": "success", "data": imageListResult.Result, "cookie": cookieJson})
			return
		} else {
			c.JSON(403, gin.H{"msg": "发生了未知错误", "data": nil, "cookie": nil})
			return
		}
	})

	r.POST("/confirmCaptcha", func(c *gin.Context) {
		var captchaBodyVo vo.CaptchaBodyVo
		_ = c.ShouldBindJSON(&captchaBodyVo)
		fmt.Println(captchaBodyVo)
		cookieStr, _ := base64.StdEncoding.DecodeString(captchaBodyVo.Cookie)
		fmt.Println(string(cookieStr))
		var cookies []*http.Cookie
		err := json.Unmarshal(cookieStr, &cookies)
		if err != nil {
			c.JSON(405, gin.H{"msg": "cookie无法解析或被篡改"})
			return
		}

		validateInfo := model.ValidateInfo{
			AccountKey: captchaBodyVo.AccountKey,
			SceneCode:  captchaBodyVo.SceneCode,
			TenantId:   captchaBodyVo.TenantId,
			UserId:     captchaBodyVo.UserId,
		}

		captchaBody := model.CaptchaBody{
			ScenesImageCode:  captchaBodyVo.ScenesImageCode,
			ScenesImageCodes: captchaBodyVo.ScenesImageCodes,
		}

		captchaResult, newCookies := utils.ConfirmCaptcha(cookies, captchaBodyVo.DeviceId, validateInfo, captchaBody)
		cookieList := make([]model.Cookie, 0)
		for _, cookie := range newCookies {
			cookieList = append(cookieList, model.Cookie{Name: cookie.Name, Value: cookie.Value})
		}
		cookieJson, _ := json.Marshal(cookieList)

		if captchaResult.Code == 200 {
			c.JSON(200, gin.H{"msg": "success", "data": captchaResult, "cookie": cookieJson})
			return
		} else {
			c.JSON(405, gin.H{"msg": "fail", "data": captchaResult, "cookie": nil})
			return
		}
	})

	r.GET("/autoSign", func(c *gin.Context) {
		cookieRaw := c.Query("cookie")
		username := c.Query("username")
		deviceId := c.Query("deviceId")
		validateTicket := c.Query("ticket")

		// 加载默认配置
		defaultConfig := config.GetBksConfig(validateTicket)

		if len(cookieRaw) == 0 {
			c.JSON(403, gin.H{"msg": "参数不全"})
			return
		}
		cookieStr, _ := base64.StdEncoding.DecodeString(cookieRaw)

		var cookies []*http.Cookie
		err := json.Unmarshal(cookieStr, &cookies)
		if err != nil {
			c.JSON(405, gin.H{"msg": "cookie无法解析或被篡改"})
			return
		}

		cookies = utils.HandleCheatModule(cookies)

		var OpJar *cookiejar.Jar
		OpJar, _ = cookiejar.New(nil)
		client := resty.NewWithClient(
			&http.Client{
				Jar: OpJar,
			},
		)
		// 获取签到任务列表
		signInfoUrl := "https://nbu.campusphere.net/wec-counselor-sign-apps/stu/sign/getStuSignInfosInOneDay"
		signInfo := vo.SignInfo{}
		_, _ = client.R().
			SetResult(&signInfo).
			SetHeaders(map[string]string{
				"User-Agent":   "Mozilla/5.0 (Macintosh; Apple Mac OS X 13_1_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.88 Safari/605.1.15",
				"Content-Type": "application/json",
			}).SetCookies(cookies).SetBody(map[string]string{}).Post(signInfoUrl)

		taskInfoForm := make(map[string]string)
		if len(signInfo.Datas.UnSignedTasks) < 1 {
			fmt.Println("没有签到任务")
			c.JSON(404, gin.H{"msg": "没有签到任务"})
			return
		} else {
			fmt.Println(signInfo.Datas.UnSignedTasks)
			taskInfoForm = map[string]string{
				"signInstanceWid": signInfo.Datas.UnSignedTasks[0].SignInstanceWid,
				"signWid":         signInfo.Datas.UnSignedTasks[0].SignWid,
				"taskName":        signInfo.Datas.UnSignedTasks[0].TaskName,
			}
		}

		// 获取签到任务详情
		taskInfoUrl := "https://nbu.campusphere.net/wec-counselor-sign-apps/stu/sign/detailSignInstance"
		taskInfo := vo.TaskInfo{}
		_, _ = client.R().
			SetResult(&taskInfo).
			SetHeaders(map[string]string{
				"User-Agent":   "Mozilla/5.0 (Macintosh; Apple Mac OS X 13_1_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.88 Safari/605.1.15",
				"Content-Type": "application/json;charset=UTF-8",
			}).SetCookies(cookies).SetBody(taskInfoForm).Post(taskInfoUrl)

		// todo 适配研究生
		if len(taskInfo.Datas.SignedStuInfo.UserId) > 9 {
			defaultConfig = config.GetYjsConfig(validateTicket)
		}

		// 解析数据
		// 绑定签到任务的 SignInstanceWid 数据
		defaultConfig.SignInstanceWid = taskInfo.Datas.SignInstanceWid

		// 绑定签到任务中附加信息的选项wid
		extraFields := taskInfo.Datas.ExtraField
		for i := 0; i < len(extraFields); i++ {
			for j := 0; j < len(extraFields[i].ExtraFieldItems); j++ {
				if defaultConfig.ExtraFieldItems[i].ExtraFieldItemValue == extraFields[i].ExtraFieldItems[j].Content {
					defaultConfig.ExtraFieldItems[i].ExtraFieldItemWid = extraFields[i].ExtraFieldItems[j].Wid
					continue
				}
			}
			if defaultConfig.ExtraFieldItems[i].ExtraFieldItemWid == 0 {
				fmt.Println("有选项找不到，检查选项内容以及选项顺序是否一致！！！")
				c.JSON(404, gin.H{"msg": "有选项找不到，可能签到问题配置需要升级！"})
				return
			}
		}

		data, _ := json.Marshal(&defaultConfig)
		fmt.Println(string(data))

		// 处理额外信息
		extensionData, _ := json.Marshal(config.GetExtensionInfo(username, deviceId, validateTicket))
		cpdailyExtension := utils.EncryptCpdailyExtension(extensionData)
		fmt.Println(string(extensionData))

		// 获取submit form
		submitForm := utils.GetSubmitForm(username, deviceId)
		submitForm.BodyString = utils.EncryptBodyString(data)
		submitFormData, _ := json.Marshal(submitForm)

		// 获取上传payload的md5摘要
		submitFormMap := make(map[string]string)
		err = json.Unmarshal(submitFormData, &submitFormMap)
		if err != nil {
			c.JSON(500, gin.H{"msg": "对提交的信息进行签名发生错误"})
			return
		}
		abstractKey := []string{"appVersion", "bodyString", "deviceId", "lat", "lon", "model", "systemName", "systemVersion", "userId"}
		payload := ""
		for i := 0; i < len(abstractKey); i++ {
			payload += url.QueryEscape(abstractKey[i]) + "=" + url.QueryEscape(submitFormMap[abstractKey[i]])
			if i < len(abstractKey)-1 {
				payload += "&"
			}
		}
		submitForm.Sign = fmt.Sprintf("%x", md5.Sum([]byte(payload)))
		submitFormMap["sign"] = submitForm.Sign
		submitFormData, _ = json.Marshal(submitForm)

		// 发送签到请求
		submitUrl := "https://nbu.campusphere.net/wec-counselor-sign-apps/stu/sign/submitSign"

		submitReq, _ := http.NewRequest("POST", submitUrl, strings.NewReader(string(submitFormData)))
		submitClient := &http.Client{}
		for _, cookie := range cookies {
			submitReq.AddCookie(cookie)
		}
		submitReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Apple Mac OS X 13_1_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.88 Safari/605.1.15")
		submitReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
		submitReq.Header.Add("Host", "https://nbu.campusphere.net/")
		submitReq.Header.Add("Cpdaily-Extension", cpdailyExtension)

		res, err := submitClient.Do(submitReq)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(res.Body)
		result, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(result))

		var autoSignResult model.AutoSignResult
		err = json.Unmarshal(result, &autoSignResult)

		if autoSignResult.Code == "0" || autoSignResult.Message == "SUCCESS" {
			c.JSON(200, gin.H{"msg": "success", "data": autoSignResult})
			return
		} else {
			c.JSON(404, gin.H{"msg": "fail", "data": autoSignResult})
		}
	})

	err := r.Run("0.0.0.0:8188")
	if err != nil {
		return
	}
}
