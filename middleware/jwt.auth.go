package middleware

import (
	"log"
	"net/http"

	"github.com/Adebayobenjamin/clean_arch/helper"
	"github.com/Adebayobenjamin/clean_arch/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthorizeJWT vaildates the user token, returns 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No auth token found", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			token, err := jwtService.ValidateToken(authHeader)
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				log.Println("Claims[user_id]: ", claims["user_id"])
				log.Panicln("Claims[issuer]: ", claims["issuer"])
			} else {
				log.Println(err)
				response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			}
		}
	}
}
