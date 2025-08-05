package domain

type SuccessResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type TimesBonusResponse struct {
	Level int   `json:"level"`
	Time  int64 `json:"time"`
}

type ChapterResponse struct {
	Chapter int `json:"chapter"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RankingResponse struct {
	NewUserRank []User `json:"newUserRank,omitempty"`
	UserRank    []User `json:"userRank,omitempty"`
	VehicleRank []User `json:"vehicleRank,omitempty"`
}

var Code_success = 200
var Code_wrong_arg = 1001
var Code_id_exist = 1002
var Code_encrypt_fail = 1003
var Code_db_error = 1004
var Code_token_error = 1005
var Code_user_not_exist = 1006
var Code_id_wrong = 1007
var Code_requirements_wrong = 1008
var Code_get_again = 1009

var Code_fail = 10001
