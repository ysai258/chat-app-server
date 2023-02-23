package utils

import (
	"fmt"
	"net/http"
	"server/internal/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to has passowrd: %w", err)
	}
	return string(hashPswd), nil
}

func CheckPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func ServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func GetTokenData(c *gin.Context) *constants.TokenClaims {
	tkn, ok := c.Get(constants.JWT_TOKEN_CLAIMS_KEY)
	if !ok {
		ServerError(c, fmt.Errorf("invalid token"))
		return nil
	}
	tokenData := tkn.(*constants.TokenClaims)
	return tokenData
}
