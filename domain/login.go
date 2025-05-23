package domain

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
	OfflineCoin  uint64 `json:"offlineCoin"`
	OfflineTime  int64  `json:"offlineTime"`
}
