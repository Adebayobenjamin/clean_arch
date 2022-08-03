package controller

import (
	"net/http"
	"strconv"

	"github.com/Adebayobenjamin/clean_arch/dto"
	"github.com/Adebayobenjamin/clean_arch/entity"
	"github.com/Adebayobenjamin/clean_arch/helper"
	"github.com/Adebayobenjamin/clean_arch/service"
	"github.com/gin-gonic/gin"
)

// Authcontroller is a contract that specifies what this controller can do
type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	// this is where you put your service
	authService service.AuthService
	jwtService  service.JWTService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please Check your credentials", "invalid credentials", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}
func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	isDuplicate := c.authService.IsDuplicateEmail(registerDTO.Email)
	if !isDuplicate {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate Email", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}

}
