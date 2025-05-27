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
