package constants

import "github.com/golang-jwt/jwt/v4"

const REQUEST_TIMEOUT_SECONDS = 2
const TOKEN_EXPIRY_HOURS = 7 * 24
const TOKEN_MAX_AGE_SECONDS = 7 * 24 * 60 * 60 // 1 week

const ROUTER_MAX_AGE_HOURS = 12

var JWT_SECRET = ""

const JWT_TOKEN_NAME = "access_token"
const JWT_TOKEN_CLAIMS_KEY = "claims"

const BASE_SERVER_DOMAIN = "localhost"
const BASE_SERVER_PORT = "8080"
const BASE_SERVER_URL = BASE_SERVER_DOMAIN + ":" + BASE_SERVER_PORT

const BASE_CLIENT_DOMAIN = "localhost"
const BASE_CLIENT_PORT = "3000"

const BASE_CLIENT_URL = BASE_CLIENT_DOMAIN + ":" + BASE_CLIENT_PORT

type TokenClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
