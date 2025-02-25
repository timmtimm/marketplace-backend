package helper

import (
	"crop_connect/util"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTCustomClaims struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var JWTSecretKey = util.GetConfig("JWT_SECRET_KEY")

func GenerateToken(uid string, role string) string {
	claims := JWTCustomClaims{
		uid,
		role,
		jwt.RegisteredClaims{
			Issuer:    "crop_connect",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWTSecretKey))
	BearerToken := "Bearer " + token
	return BearerToken
}

func GetPayloadToken(token string) (JWTCustomClaims, error) {
	claims := JWTCustomClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return JWTCustomClaims{}, err
	}

	if !tkn.Valid {
		return JWTCustomClaims{}, errors.New("invalid token")
	}

	return claims, nil
}

func GetPayloadFromToken(c echo.Context) (JWTCustomClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	claims, err := GetPayloadToken(token)
	if err != nil {
		return JWTCustomClaims{}, err
	}

	return claims, nil
}

func GetUIDFromToken(c echo.Context) (primitive.ObjectID, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	claims, err := GetPayloadToken(token)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, err := primitive.ObjectIDFromHex(claims.UID)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user id")
	}

	return id, nil
}
