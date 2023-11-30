package service

type UserSignupRequest struct {
	UserAccount     string `json:"userAccount" binding:"required"`
	UserPassword    string `json:"userPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=UserPassword"`
}

type UserLoginRequest struct {
	UserAccount string `json:"userAccount" binding:"required"`
	Password    string `json:"userPassword" binding:"required"`
}
