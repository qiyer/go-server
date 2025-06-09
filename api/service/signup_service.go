package service

import (
	"net/http"
	"time"

	"go-server/bootstrap"
	"go-server/domain"
	"go-server/redis"
	"go-server/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var Env *bootstrap.Env

func Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_wrong_arg,
			Message: "请求参数错误",
		})
		return
	}

	_, err = repository.GetByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_id_exist,
			Message: "邮箱已经存在",
		})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_encrypt_fail,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	request.Password = string(encryptedPassword)

	account := domain.Account{
		ID:        primitive.NewObjectID(),
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: time.Now(),
	}

	err = repository.CreateAccount(c, &account)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	user := domain.User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      request.Name,
		Level:     1,
		Build: domain.Build{
			Level:        1,
			DisplayLevel: 1,
		},
		Vehicle: domain.Vehicle{
			Level:        1,
			DisplayLevel: 1,
		},
		Girls:                domain.InitGirls,
		QuickEarn:            1, // 快速收益 默认1级
		ContinuousClick:      1, // 连续点击 默认1次
		TimesBonus:           1, // 1倍收益 默认1倍
		TimesBonusTimeStamp:  0,
		Coins:                0,
		ConsecutiveLoginDays: 0,
		LastLoginDate:        "2006-01-02", // 设置为当前日期
		LastClickTimeStamp:   time.Now().Unix(),
	}

	err = repository.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	userMapping := domain.UserMapping{
		PlatformID: account.ID.Hex(),
		UserId:     user.ID,
		Platform:   "email",
		CreateAt:   user.CreatedAt,
	}

	err = repository.CreateUserMapping(c, &userMapping)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_db_error,
			Message: "系统错误，请稍后重试",
		})
		return
	}

	accessToken, err := repository.CreateAccessToken(&user, Env.AccessTokenSecret, Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			Code:    domain.Code_token_error,
			Message: "生成AccessToken失败:",
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

	signupResponse := domain.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	repository.SetLastLoginCache(user.ID.Hex(), time.Now().Unix())
	repository.SetUserCache(user.ID.Hex(), user)

	redis.CacheUserData(&user)

	c.JSON(http.StatusOK, domain.Response{
		Code: domain.Code_success,
		Data: signupResponse,
	})
}
