package user

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/constants"
	"server/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
	timeout time.Duration
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
		timeout: time.Duration(constants.REQUEST_TIMEOUT_SECONDS) * time.Second,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	userNameExist, err := h.Service.CheckUserName(ctx, req.Username)
	if err != nil {
		utils.ServerError(c, err)
		return
	}
	if userNameExist {
		utils.BadRequest(c, fmt.Errorf("username already exist"))
		return
	}
	emailExist, err := h.Service.CheckEmail(ctx, req.Email)
	if err != nil {
		utils.ServerError(c, err)
		return
	}
	if emailExist {
		utils.BadRequest(c, fmt.Errorf("email already exist"))
		return
	}
	res, err := h.Service.CreateUser(ctx, &req)
	if err != nil {
		utils.ServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	emailExist, err := h.Service.CheckEmail(ctx, req.Email)
	if err != nil {
		utils.ServerError(c, err)
		return
	}
	if !emailExist {
		utils.BadRequest(c, fmt.Errorf("invalid email"))
		return
	}

	res, err := h.Service.Login(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid Password") {
			utils.BadRequest(c, err)
			return
		}
		utils.ServerError(c, err)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constants.JWT_TOKEN_NAME, res.accessToken, constants.TOKEN_MAX_AGE_SECONDS, "/", constants.BASE_CLIENT_URL, true, true)
	c.JSON(http.StatusOK, LoginResponse{ID: res.ID, Username: res.Username})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(constants.JWT_TOKEN_NAME, "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *Handler) GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()
	tokenData := utils.GetTokenData(c)
	if tokenData == nil {
		return
	}

	user, err := h.Service.GetUser(ctx, tokenData.ID)
	if err != nil {
		utils.ServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
