package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gin-dbo/model/login"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	BEARER_SCHEMA     = "Bearer"
	Authorization     = "Authorization"
	ErrorInvalidToken = "token is not valid"
	ErrorMissingAuth  = "missing authorization"
	MethodGet         = "GET"
	MethodPost        = "POST"
	MethodDelete      = "DELETE"
	JwtSecretKey      = "JWT_SECRET_KEY"
	JwtIssuer         = "JWT_ISSUER"
	JwtClaims         = "JWT_CLAIMS"
)

type JWTService interface {
	GenerateToken(p *login.User) string
	ValidateToken(token string) (*AuthCustomClaims, error)
}
type AuthCustomClaims struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	CustomerId string `json:"customer_id"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

type Response struct {
	Code    int    `json:"code" example:"401"`
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"unauthorized"`
}

func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: os.Getenv(JwtSecretKey),
		issuer:    os.Getenv(JwtIssuer),
	}
}

func (service *jwtServices) GenerateToken(p *login.User) string {
	claims := &AuthCustomClaims{
		Username:   p.Username,
		Role:       p.Role,
		CustomerId: p.CustomerId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Subject:   p.Username,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return err.Error()
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*AuthCustomClaims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &AuthCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(JwtSecretKey)), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token : %v", err)
	}
	jwtClaims, ok := token.Claims.(*AuthCustomClaims)
	if ok && token.Valid {
		return jwtClaims, err
	} else {
		return nil, fmt.Errorf(ErrorInvalidToken)
	}
}

func AuthorizeJWT() func(*gin.Context) {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.GetHeader(Authorization), " ")
		if authHeader[0] != BEARER_SCHEMA {
			err := errors.New(ErrorMissingAuth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{Code: http.StatusUnauthorized, Success: false, Message: err.Error()})
			return
		} else {
			jwtSegment := strings.Split(authHeader[1], ".")
			if authHeader[1] != "" && len(jwtSegment) == 3 {
				tokenString := authHeader[1]
				claims, err := JWTAuthService().ValidateToken(tokenString)
				if err == nil {
					c.Set(JwtClaims, claims)
					claims.validatePath(c)
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{Code: http.StatusUnauthorized, Success: false, Message: err.Error()})
					return
				}
			} else {
				err := errors.New(ErrorMissingAuth)
				c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{Code: http.StatusUnauthorized, Success: false, Message: err.Error()})
				return
			}
		}
	}
}

func (claims *AuthCustomClaims) validatePath(c *gin.Context) {
	if claims.Role == "admin" {
		c.Next()
	} else {
		if strings.Contains(c.Request.URL.Path, "order") {
			c.Next()
		} else if strings.Contains(c.Request.URL.Path, "customer") && c.Request.Method != MethodDelete && c.Request.Method != MethodPost {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{Code: http.StatusUnauthorized, Success: false, Message: "this user doesn't have access to this endpoint"})
			return
		}
	}
}
