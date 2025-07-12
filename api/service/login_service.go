package service

import (
	"go-server/domain"
	"go-server/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	user, account, err := repository.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	repository.SetLastLoginCache(user.ID.Hex(), time.Now().Unix())
	repository.SetUserCache(user.ID.Hex(), user)

	// redis.CacheUserData(&user)

	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_encrypt_fail,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	now := time.Now().Format("2006-01-02")

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

	accessToken, err := repository.CreateAccessToken(&user, Env.AccessTokenSecret, Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_token_error,
			Message: "生成AccessToken失败",
		})
		return
	}

	refreshToken, err := repository.CreateRefreshToken(&user, Env.RefreshTokenSecret, Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_token_error,
			Message: "生成refreshToken失败: ",
		})
		return
	}

	var offlineTime = time.Now().Unix() - user.UpdatedAt.Unix()
	var secCoin = domain.GetSecCoin(user)
	var offlineCoin = domain.GetOfflineCoin(secCoin, uint64(offlineTime))

	user = domain.GetNewUser(user)

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		OfflineCoin:  offlineCoin,
		OfflineTime:  offlineTime,
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: loginResponse,
	})

}

func RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	id, err := repository.ExtractIDFromToken(request.RefreshToken, Env.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	user, err := repository.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_user_not_exist,
			Message: "用户不存在",
		})
		return
	}

	accessToken, err := repository.CreateAccessToken(&user, Env.AccessTokenSecret, Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_token_error,
			Message: "生成AccessToken失败: ",
		})
		return
	}

	refreshToken, err := repository.CreateRefreshToken(&user, Env.RefreshTokenSecret, Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_token_error,
			Message: "生成refreshToken失败: ",
		})
		return
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: refreshTokenResponse,
	})
}

func UpdateUserInfo(c *gin.Context) {
	var grow domain.UserInfoRequest

	err := c.ShouldBind(&grow)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// 将字符串转换为primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(grow.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := repository.UpdateUserCoins(c, userID, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	repository.SetUserCache(user.ID.Hex(), user)

	c.JSON(http.StatusOK, user)
}
