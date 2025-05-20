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
	var grow domain.CoinAutoQueueRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&grow)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.UpdateUserCoins(c, userID, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetOfflineCoin(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User not found"})
		return
	}

	var offlineTime = time.Now().Unix() - user.UpdatedAt.Unix()
	var secCoin = domain.GetSecCoin(user)
	var offlineCoin = domain.GetOfflineCoin(secCoin, uint64(offlineTime))

	nuser, err := repository.UpdateUserCoins(c, userID, offlineCoin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nuser)
}

func CheckIn(c *gin.Context) {
	user_id := c.GetString("x-user-id")
	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.CheckIn(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ClaimOnlineRewards(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.ClaimOnlineRewards(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func LevelUp(c *gin.Context) {
	var res domain.LevelUpRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User not found"})
		return
	}

	var costCoin uint64 = 1

	//主角升级
	if res.RoleID == 10000 {
		costCoin = domain.RoleLevelCost(user.Level)
		if user.Coins < costCoin {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Not enough coins"})
			return
		}
		nuser, err := repository.LevelUp(c, userID, res.Level, costCoin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, nuser)
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
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "秘书已经满级"})
			return
		}

		updatedGirls := strings.Join(newGirls, "")
		if level == -1 {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "秘书不存在"})
			return
		}

		costCoin = domain.GirlLevelCost(res.RoleID, level)
		if user.Coins < costCoin {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "girl:Not enough coins"})
			return
		}

		nuser, err := repository.RoleLevelUp(c, userID, updatedGirls, costCoin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, nuser)
	}
}

func PassChapter(c *gin.Context) {
	var res domain.PassChapterRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.PassChapter(c, userID, res.Chapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Ranking(c *gin.Context) {

	users, err := repository.Ranking(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateTask(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")
	task.ID = primitive.NewObjectID()

	task.UserID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err = repository.CreateTask(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
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
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err = repository.UpgradeApartment(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Apartment upgraded successfully", Code: 200,
	})
}

func UnLockRole(c *gin.Context) {
	var res domain.UnLockRoleRequest
	user_id := c.GetString("x-user-id")
	err := c.ShouldBind(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.GetByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User not found"})
		return
	}

	var newGirl = fmt.Sprintf("%d", res.RoleID)
	if strings.Contains(user.Girls, newGirl) {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "秘书已存在"})
		return
	}

	var isInConfig = false
	for _, girl := range domain.Girls {
		if girl.GirlId == res.RoleID {
			isInConfig = true
			break
		}
	}

	if !isInConfig {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "配置表格里不存在该秘书"})
		return
	}

	if !domain.GirlUnlockCheckNeeds(res.RoleID, user) {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "主角或相关角色等级不足"})
		return
	}

	var updatedGirls = fmt.Sprintf("%s%d:0", user.Girls, res.RoleID)

	nuser, err := repository.RoleLevelUp(c, userID, updatedGirls, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nuser)
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
