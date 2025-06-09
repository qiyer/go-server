package service

import (
	"fmt"
	"go-server/domain"
	"go-server/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CoinAutoGrowing(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	now := time.Now().Format("2006-01-02")

	fmt.Printf("Login time %s\n", now)

	if user.LastLoginDate != now {
		day := user.ConsecutiveLoginDays + 1
		days := user.Days
		if len(user.Days) > 7 {

		} else {
			days = domain.CheckInDays(user, day)
		}
		user.Days = days
		_ = repository.UpdateUserDays(c, user.ID, days, day, now)
	}

	timeNow := time.Now().Unix() - repository.GetLastLoginCache(user_id)

	var onlineTime = user.OnlineTime + int(timeNow)

	var secCoin = domain.GetSecCoin(user)
	var index = uint64(1)
	var bonusTime = int64(0)
	lastUpdateStamp := user.UpdatedAt.Unix()
	if lastUpdateStamp < user.TimesBonusTimeStamp {
		var timeDiff = user.TimesBonusTimeStamp - lastUpdateStamp
		if timeDiff > 5 {
			bonusTime = user.TimesBonusTimeStamp - lastUpdateStamp - 5
		} else {
			bonusTime = user.TimesBonusTimeStamp - lastUpdateStamp - timeDiff
		}
		index = uint64(user.TimesBonus)
		bonusTime = bonusTime + time.Now().Unix()
	}
	// 多倍收益计算需要传入
	addCoin := domain.GetOnlineCoin(secCoin, uint64(onlineTime), index)

	nuser, err := repository.UpdateUserCoinsWithTime(c, userID, uint64(addCoin), string(onlineTime), bonusTime)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	repository.SetLastLoginCache(user_id, time.Now().Unix())

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func ClickEarn(c *gin.Context) {
	var res domain.ClickEarnRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var timediff = user.LastClickTimeStamp - user.UpdatedAt.Unix()
	// 假定点击时间是 2s 一次 ， auto coin growing 5s 一次
	if int(timediff) > 3 || int(timediff) < -3 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "您涉嫌作弊，点击时间不符合要求",
		})
		return
	}
	// 假定每秒最多点击 5 次
	if int(res.Clicker) > 10 || int(res.Clicker) < 1 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "您涉嫌作弊，点击次数不符合要求",
		})
		return
	}

	var index = uint64(1)
	lastUpdateStamp := user.UpdatedAt.Unix()
	if lastUpdateStamp > user.TimesBonusTimeStamp {
		index = uint64(user.TimesBonus)
	}

	addCoin := domain.GetClickCoin(user, 1, uint64(res.Clicker), index)
	timenow := time.Now().Unix()
	nuser, err := repository.UpdateUserCoinsWithClick(c, userID, uint64(addCoin), timenow)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func GetOfflineCoin(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var offlineTime = time.Now().Unix() - user.UpdatedAt.Unix()
	var secCoin = domain.GetSecCoin(user)
	var offlineCoin = domain.GetOfflineCoin(secCoin, uint64(offlineTime))

	nuser, err := repository.UpdateUserCoins(c, userID, offlineCoin)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func CheckIn(c *gin.Context) {
	var res domain.CheckInRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	if res.Id < 1 || res.Id > 7 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	dayStr := fmt.Sprintf("%d", res.Id)

	isCheck, days := domain.CheckIn(user, dayStr)

	if !isCheck {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "已签到或是条件不符合",
		})
		return
	}

	err = repository.UpdateDays(c, userID, days)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	var index = len(user.Days)
	var reward = domain.DayBonuses[index]

	if reward.Type == "coin" {

		nuser, err := repository.UpdateUserCoins(c, userID, uint64(reward.Bonus))
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]any{"bonus_type": "coin", "user": nuser},
		})
		return
	} else if reward.Type == "level" {
		// 升级奖励

		nuser, err := repository.LevelUp(c, userID, int(reward.Bonus), 0)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]any{"bonus_type": "level", "user": nuser},
		})
		return
	} else if reward.Type == "role" {
		var newGirl = fmt.Sprintf("%s:", reward.Bonus)
		for _, girl := range user.Girls {
			if strings.Contains(girl, newGirl) {
				c.JSON(http.StatusOK, domain.Response{
					Code:    domain.Code_requirements_wrong,
					Message: "角色已解锁,领取失败",
				})
				return
			}
		}
		var updatedGirls = append(user.Girls, fmt.Sprintf("%d:0", reward.Bonus))
		nuser, err := repository.RoleLevelUp(c, userID, updatedGirls, 0)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]any{"bonus_type": "role", "user": nuser},
		})
		return
	} else if reward.Type == "box" {
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]string{"bonus_type": "box", "num": "5"},
		})
		return
	} else if reward.Type == "click" {
		level, _ := repository.ContinuousClick(c, userID, 1)
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]string{"bonus_type": "click", "level": string(level)},
		})
		return
	}

	c.JSON(http.StatusOK, "签到成功，获得奖励！")
}

func ClaimOnlineRewards(c *gin.Context) {
	var res domain.OnlineRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	if res.Id < 1 || res.Id > 5 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	status := user.OnlineRewards[res.Id-1]

	if status == 1 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_get_again,
			Message: "已领取该奖励",
		})
		return
	}

	var reward = domain.OnlineBonuses[res.Id-1]

	if user.OnlineTime < reward.Min {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "条件不符合，在线时长不足",
		})
		return
	}

	user.OnlineRewards[res.Id-1] = 1 // 更新状态为已领取

	err = repository.UpdateOnlineRewards(c, userID, user.OnlineRewards)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	if reward.Type == "coin" {

		nuser, err := repository.UpdateUserCoins(c, userID, uint64(reward.Bonus))
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, nuser)
	} else if reward.Type == "level" {
		// 升级奖励

		nuser, err := repository.LevelUp(c, userID, int(reward.Bonus), 0)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, nuser)
		return
	} else if reward.Type == "role" {
		var newGirl = fmt.Sprintf("%s:", reward.Bonus)

		for _, girl := range user.Girls {
			if strings.Contains(girl, newGirl) {
				c.JSON(http.StatusOK, domain.Response{
					Code:    domain.Code_db_error,
					Message: "角色已解锁,领取失败",
				})
				return
			}
		}
		var updatedGirls = append(user.Girls, fmt.Sprintf("%d:0", reward.Bonus))
		nuser, err := repository.RoleLevelUp(c, userID, updatedGirls, 0)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}
		c.JSON(http.StatusOK, nuser)
		return
	} else if reward.Type == "box" {
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]string{"bonus_type": "box", "num": "5"},
		})
		return
	} else if reward.Type == "click" {
		level, _ := repository.ContinuousClick(c, userID, 1)
		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: map[string]string{"bonus_type": "click", "level": string(level)},
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func LevelUp(c *gin.Context) {
	var res domain.LevelUpRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var costCoin uint64 = 1

	//主角升级
	if res.RoleID == 10000 {
		costCoin = domain.RoleLevelCost(user.Level)
		if user.Coins < costCoin {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_requirements_wrong,
				Message: "金币不足",
			})
			return
		}
		nuser, err := repository.LevelUp(c, userID, res.Level, costCoin)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}

		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: nuser,
		})
		return
	} else {
		//秘书升级
		var level = -1
		girls := domain.ParseGirls(user.Girls)
		var newGirls []string
		for _, girl := range girls {
			if girl.GirlId == res.RoleID {
				level = int(girl.Level)
				newGirls = append(newGirls, fmt.Sprintf("%d:%d,", girl.GirlId, girl.Level+1))
			} else {
				newGirls = append(newGirls, fmt.Sprintf("%d:%d,", girl.GirlId, girl.Level))
			}
		}

		if level >= 100 {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "秘书已经满级",
			})
			return
		}

		if level == -1 {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "秘书不存在",
			})
			return
		}

		costCoin = domain.GirlLevelCost(res.RoleID, level)
		if user.Coins < costCoin {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_requirements_wrong,
				Message: "金币不足",
			})
			return
		}

		nuser, err := repository.RoleLevelUp(c, userID, newGirls, costCoin)
		if err != nil {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_db_error,
				Message: "系统错误，请稍后重试",
			})
			return
		}

		c.JSON(http.StatusOK, domain.Response{
			Code: domain.Code_success,
			Data: nuser,
		})
	}
}

func PassChapter(c *gin.Context) {
	var res domain.PassChapterRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	chapter, err := repository.PassChapter(c, userID, res.Chapter)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	var response = domain.ChapterResponse{}
	response.Chapter = chapter

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: response,
	})
}

func Ranking(c *gin.Context) {

	users, err := repository.Ranking(c)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: users,
	})
}

func CaiShen(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}
	coin, err := repository.CaiShen(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "成功获得:" + fmt.Sprintf("%d", coin),
		Data:    user,
	})
}

func TimesBonus(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}
	resp, err := repository.TimesBonus(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "成功:" + fmt.Sprintf("%d", resp.Level),
		Data:    user,
	})
}

func ContinuousClick(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}
	level, err := repository.ContinuousClick(c, userID, 1)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "成功升级:" + fmt.Sprintf("%d", level),
		Data:    user,
	})
}

func QuickEarn(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}
	coin, err := repository.QuickEarn(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "成功获得:" + fmt.Sprintf("%d", coin),
		Data:    user,
	})

}

func CreateTask(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	userID := c.GetString("x-user-id")
	task.ID = primitive.NewObjectID()

	task.UserID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	err = repository.CreateTask(c, &task)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
}

func UpgradeApartment(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	err = repository.UpgradeApartment(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "小区升级成功",
		Data:    user,
	})
}

func ChangeVehicle(c *gin.Context) {
	var res domain.VehicleDisplayRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	err = repository.ChangeVehicleVehicle(c, userID, res.DisplayLevel)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code:    domain.Code_success,
		Message: "坐骑升级成功",
		Data:    user,
	})
}

func UpgradeVehicle(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	err = repository.UpgradeVehicle(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: user,
	})
}

func UnLockRole(c *gin.Context) {
	var res domain.UnLockRoleRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var newGirl = fmt.Sprintf("%d", res.RoleID)
	for _, girl := range user.Girls {
		if strings.Contains(girl, newGirl) {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_requirements_wrong,
				Message: "秘书已存在",
			})
			return
		}
	}

	var isInConfig = false
	for _, girl := range domain.Girls {
		if girl.GirlId == res.RoleID {
			isInConfig = true
			break
		}
	}

	if !isInConfig {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "不存在该秘书",
		})
		return
	}

	if !domain.GirlUnlockCheckNeeds(res.RoleID, user) {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "主角或相关角色等级不足",
		})
		return
	}

	var updatedGirls = append(user.Girls, fmt.Sprintf("%d:0", res.RoleID))

	nuser, err := repository.RoleLevelUp(c, userID, updatedGirls, 0)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func UnLockVehicle(c *gin.Context) {
	// 	var res domain.UnLockVehicleRequest
	// 	user_id := c.GetString("x-user-id")
	// 	err := c.ShouldBind(&res)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	// 		return
	// 	}

	// 	// 将字符串转换为primitive.ObjectID
	// 	userID, err := primitive.ObjectIDFromHex(user_id)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	// 		return
	// 	}

	// 	user, err := repository.GetByID(c, userID)
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User not found"})
	// 		return
	// 	}

	// 	var newVehicle = fmt.Sprintf("%d", res.VehicleID)
	// 	if strings.Contains(user.Vehicles, newVehicle) {
	// 		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "坐骑已存在"})
	// 		return
	// 	}

	// 	var isInConfig = false
	// 	for _, vehicle := range domain.Vehicles {
	// 		if vehicle.ID == res.VehicleID {
	// 			isInConfig = true
	// 			break
	// 		}
	// 	}

	// 	if !isInConfig {
	// 		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "配置表格里不存在该坐骑"})
	// 		return
	// 	}

	// 	success, coin := domain.VehicleUnlockCheckNeeds(res.VehicleID, user)

	// 	if !success {
	// 		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "角色等级不足或是金币不足"})
	// 		return
	// 	}

	// 	var updatedVehicles = fmt.Sprintf("%s;%d", user.Vehicles, res.VehicleID)

	// 	nuser, err := repository.UnLockVehicle(c, userID, updatedVehicles, coin)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	// 		return
	// 	}

	// c.JSON(http.StatusOK, nuser)
}

func UnLockCapital(c *gin.Context) {
	var res domain.UnLockCapitalRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var newCapital = fmt.Sprintf("%d:%d", res.CapitalID, time.Now().Unix())
	var checkCapital = fmt.Sprintf("%d:", res.CapitalID)
	for _, capital := range user.Capitals {
		if strings.Contains(capital, checkCapital) {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_requirements_wrong,
				Message: "资产已存在",
			})
			return
		}
	}

	var isInConfig = false
	for _, capital := range domain.Capitals {
		if capital.ID == res.CapitalID {
			isInConfig = true
			break
		}
	}

	if !isInConfig {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "不存在该资产",
		})
		return
	}

	success, coin := domain.CapitalUnlockCheckNeeds(res.CapitalID, user)

	if !success {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "金币不足",
		})
		return
	}

	var updatedCapitals = append(user.Capitals, newCapital)

	nuser, err := repository.UnLockCapital(c, userID, updatedCapitals, coin)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func GetCapitalIncome(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	if len(user.Capitals) == 0 {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_requirements_wrong,
			Message: "资产不存在",
		})
		return
	}

	coin, caps := domain.GetCapitalIncome(user)

	nuser, err := repository.SellCapital(c, userID, caps, coin)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})
}

func SellCapital(c *gin.Context) {
	var res domain.UnLockCapitalRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_wrong,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	var checkCapital = fmt.Sprintf("%d:", res.CapitalID)

	for _, capital := range user.Capitals {
		if !strings.Contains(capital, checkCapital) {
			c.JSON(http.StatusOK, domain.Response{
				Code:    domain.Code_requirements_wrong,
				Message: "资产不存在",
			})
			return
		}
	}

	coin, caps := domain.SellCapital(res.CapitalID, user)

	nuser, err := repository.SellCapital(c, userID, caps, coin)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: nuser,
	})

}

func FetchTask(c *gin.Context) {
	userID := c.GetString("x-user-id")

	tasks, err := repository.FetchTaskByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
