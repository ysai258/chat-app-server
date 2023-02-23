package user

import (
	"context"
	"fmt"
	"server/internal/constants"
	"server/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (h *Handler) ValidateJWT(tokenString string) (*constants.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &constants.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return signing key
		return []byte(constants.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*constants.TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *service) CheckUserName(ctx context.Context, userName string) (bool, error) {
	userExist, err := s.Repository.UserNameExist(ctx, userName)
	if err != nil {
		return false, err
	}
	return userExist, nil
}

func (s *service) CheckEmail(ctx context.Context, email string) (bool, error) {
	userExist, err := s.Repository.EmailExist(ctx, email)
	if err != nil {
		return false, err
	}
	return userExist, nil
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	r, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	res := &CreateUserResponse{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}
	return res, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid Password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, constants.TokenClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.TOKEN_EXPIRY_HOURS * time.Hour)),
		},
	})

	signedTokem, err := token.SignedString([]byte(constants.JWT_SECRET))
	if err != nil {
		return nil, err
	}
	res := &LoginResponse{}
	res.ID = strconv.Itoa(int(user.ID))
	res.Username = user.Username
	res.accessToken = signedTokem
	return res, nil
}

func (s *service) GetUser(ctx context.Context, id int64) (*CreateUserResponse, error) {
	user, err := s.Repository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{ID: strconv.Itoa(int(user.ID)), Email: user.Email, Username: user.Username}, nil
}
