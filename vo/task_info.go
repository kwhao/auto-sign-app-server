package vo

import "autosign/model"

type TaskInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		SignInstanceWid string `json:"signInstanceWid"`
		SignMode        int    `json:"signMode"`
		SignRate        string `json:"signRate"`
		SignCondition   int    `json:"signCondition"`
		TaskType        string `json:"taskType"`
		TaskName        string `json:"taskName"`
		TaskDesc        string `json:"taskDesc"`
		QrCodeRcvdUsers []struct {
			TargetWid      string `json:"targetWid"`
			TargetType     string `json:"targetType"`
			TargetName     string `json:"targetName"`
			TargetGrade    string `json:"targetGrade"`
			TargetDegree   string `json:"targetDegree"`
			TargetUserType string `json:"targetUserType"`
			RoleWid        string `json:"roleWid"`
		} `json:"qrCodeRcvdUsers"`
		SenderUserName      string `json:"senderUserName"`
		CurrentTime         string `json:"currentTime"`
		SingleTaskBeginTime string `json:"singleTaskBeginTime"`
		SingleTaskEndTime   string `json:"singleTaskEndTime"`
		RateSignDate        string `json:"rateSignDate"`
		RateTaskBeginTime   string `json:"rateTaskBeginTime"`
		RateTaskEndTime     string `json:"rateTaskEndTime"`
		SignStatus          string `json:"signStatus"`
		SignTime            string `json:"signTime"`
		SignPhotoUrl        string `json:"signPhotoUrl"`
		SignType            string `json:"signType"`
		ChangeTime          string `json:"changeTime"`
		ChangeActorName     string `json:"changeActorName"`
		SignPlaceSelected   []struct {
			Wid            string `json:"wid"`
			PlaceWid       string `json:"placeWid"`
			Address        string `json:"address"`
			Longitude      string `json:"longitude"`
			Latitude       string `json:"latitude"`
			Radius         int    `json:"radius"`
			CreatorUserWid string `json:"creatorUserWid"`
			CreatorUserId  string `json:"creatorUserId"`
			CreatorName    string `json:"creatorName"`
			CurrentStatus  string `json:"currentStatus"`
			IsShare        string `json:"isShare"`
		} `json:"signPlaceSelected"`
		IsPhoto       int    `json:"isPhoto"`
		Photograph    string `json:"photograph"`
		DownloadUrl   string `json:"downloadUrl"`
		LeaveAppUrl   string `json:"leaveAppUrl"`
		CatQrUrl      string `json:"catQrUrl"`
		SignAddress   string `json:"signAddress"`
		Longitude     string `json:"longitude"`
		Latitude      string `json:"latitude"`
		IsMalposition int    `json:"isMalposition"`
		SignedStuInfo struct {
			UserWid           string                   `json:"userWid"`
			UserId            string                   `json:"userId"`
			UserName          string                   `json:"userName"`
			Sex               string                   `json:"sex"`
			Nation            string                   `json:"nation"`
			Mobile            string                   `json:"mobile"`
			Grade             string                   `json:"grade"`
			Dept              string                   `json:"dept"`
			Major             string                   `json:"major"`
			Cls               string                   `json:"cls"`
			SchoolStatus      string                   `json:"schoolStatus"`
			Malposition       string                   `json:"malposition"`
			StuDormitoryVo    model.StuDormitoryVo     `json:"stuDormitoryVo"`
			ExtraFieldItemVos []model.ExtraFieldItemVo `json:"extraFieldItemVos"`
		} `json:"signedStuInfo"`
		IsNeedExtra       int                      `json:"isNeedExtra"`
		IsAllowUpdate     bool                     `json:"isAllowUpdate"`
		UpdateLimit       int                      `json:"updateLimit"`
		LeftNum           int                      `json:"leftNum"`
		ExtraField        []model.ExtraField       `json:"extraField"`
		ExtraFieldItemVos []model.ExtraFieldItemVo `json:"extraFieldItemVos"`
	} `json:"datas"`
}
