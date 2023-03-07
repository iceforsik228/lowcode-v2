package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lowcode-v2/initializers"
	"lowcode-v2/models"
	"lowcode-v2/utils"
	"net/http"
	"time"
)

func RequireAuth(context *gin.Context) {
	// get the cookie off req
	tokenString, err := context.Cookie("Authorization")
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	// parse token (decode -> validate -> verify)
	token, err := utils.ParseToken(tokenString)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		// attach to req
		context.Set("user", user)

		// continue
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
}
