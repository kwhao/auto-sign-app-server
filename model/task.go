package model

type Task struct {
	StuSignWid          string `json:"stuSignWid"`
	SignInstanceWid     string `json:"signInstanceWid"`
	SignWid             string `json:"signWid"`
	SignRate            string `json:"signRate"`
	TaskType            string `json:"taskType"`
	TaskName            string `json:"taskName"`
	SenderUserName      string `json:"senderUserName"`
	SignStatus          string `json:"signStatus"`
	IsMalposition       string `json:"isMalposition"`
	IsLeave             string `json:"isLeave"`
	LeMobileUrl         string `json:"leMobileUrl"`
	CurrentTime         string `json:"currentTime"`
	SingleTaskBeginTime string `json:"singleTaskBeginTime"`
	SingleTaskEndTime   string `json:"singleTaskEndTime"`
	RateSignDate        string `json:"rateSignDate"`
	RateTaskBeginTime   string `json:"rateTaskBeginTime"`
	RateTaskEndme       string `json:"rateTaskEndme"`
}
