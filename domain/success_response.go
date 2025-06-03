package domain

type SuccessResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type TimesBonusResponse struct {
	Level int `json:"level"`
	Time  int `json:"time"`
}

type ChapterResponse struct {
	Chapter int `json:"chapter"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var Code_success = 200
var Code_fail = 10001
