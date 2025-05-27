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
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = repository.GetByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "Account already exists with the given email"})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
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
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user := domain.User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Level:     1,
		Build: domain.Build{
			Level:        1,
			DisplayLevel: 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Girls:             domain.InitGirls,
		QuickEarn:         1, // 快速收益 默认1级
		ContinuousClick:   1, // 连续点击 默认1次
		TimesBonus:        1, // 1倍收益 默认1倍
		TimesBonusSeconds: 0,
		Coins:             0,
	}

	err = repository.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
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
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := repository.CreateAccessToken(&user, Env.AccessTokenSecret, Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := repository.CreateRefreshToken(&user, Env.RefreshTokenSecret, Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
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

	c.JSON(http.StatusOK, signupResponse)
}
