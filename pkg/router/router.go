package router

import (
	"fmt"
	"net/http"
	"server/internal/constants"
	"server/internal/user"
	"server/internal/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	r = gin.Default()

	// Set up CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://%v:%v", constants.BASE_CLIENT_DOMAIN, constants.BASE_CLIENT_PORT)},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == fmt.Sprintf("http://%v:%v", constants.BASE_CLIENT_DOMAIN, constants.BASE_CLIENT_PORT)
		},
		MaxAge: constants.ROUTER_MAX_AGE_HOURS * time.Hour,
	}))

	// Set up JWT middleware to validate token
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/signup" {
			return // skip middleware for login and signup requests
		}

		// extract JWT from request headers or cookie
		tokenString := ""
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			tokenString = authHeader
		} else if cookie, err := c.Request.Cookie(constants.JWT_TOKEN_NAME); err == nil {
			tokenString = cookie.Value
		}

		// validate JWT
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		claims, err := userHandler.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// add claims to context for other handlers to use
		c.Set(constants.JWT_TOKEN_CLAIMS_KEY, claims)
	})

	// Set up routes
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	r.GET("/logout", userHandler.Logout)

	r.GET("/loggedUser", userHandler.GetUser)
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)

}

func Start(addr string) error {
	return r.Run(addr)
}
