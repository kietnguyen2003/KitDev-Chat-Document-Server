package interfaces

import (
	"server/internal/application/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		errorResponse(c, 400, "Can not bind request")
		return
	}

	resRegister, err := ah.authService.Register(c.Request.Context(), *reqToDto(req))
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}

	// resStorage, err :=
	successResponse(c, 200, "Register success", resRegister)
}

func (ah *AuthHandler) SignIn(c *gin.Context) {
	var req signinRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		errorResponse(c, 400, "Can not bind request")
		return
	}

	res, err := ah.authService.SignIn(c.Request.Context(), auth.SignInRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	successResponse(c, 200, "Register success", res)
}

func (ah *AuthHandler) GetMe(c *gin.Context) {

}

func reqToDto(req registerRequest) *auth.RegisterRequest {
	return &auth.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
	}
}
