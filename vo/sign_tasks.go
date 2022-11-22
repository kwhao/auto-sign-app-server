package vo

import "autosign/model"

type SignInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		DayInMonth    string       `json:"dayInMonth"`
		CodeRcvdTasks []model.Task `json:"codeRcvdTasks"`
		SignedTasks   []model.Task `json:"signedTasks"`
		UnSignedTasks []model.Task `json:"unSignedTasks"`
		LeaveTasks    []model.Task `json:"leaveTasks"`
	} `json:"datas"`
}
