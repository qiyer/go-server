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

type RankResponse struct {
	Ranks     RankingResponse `json:"ranks,omitempty"`
	UserRK    int             `json:"userRK,omitempty"`    // 用户在排行榜中的排名
	NewUserRK int             `json:"newUserRK,omitempty"` // 新用户排行榜中的排名
	VehicleRK int             `json:"vehicleRK,omitempty"` // 坐骑排行榜中的排名

	IsUserRK    bool `json:"isUserRK,omitempty"`    // 是否在用户排行榜中
	IsNewUserRK bool `json:"isNewUserRK,omitempty"` // 是否在新用户排行榜中
	IsVehicleRK bool `json:"isVehicleRK,omitempty"` // 是否在

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
