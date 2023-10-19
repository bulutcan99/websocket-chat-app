package token

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"
)

var (
	JWT_SECRET_KEY               = &env.Env.JwtSecretKey
	JWT_SECRET_KEY_EXPIRE_MINUTE = &env.Env.JwtSecretKeyExpireMinutesCount
	JWT_SECRET_KEY_EXPIRE_HOUR   = &env.Env.JwtRefreshKeyExpireHoursCount
	privateKey                   = []byte(*JWT_SECRET_KEY)
	minutesCount                 = *JWT_SECRET_KEY_EXPIRE_MINUTE
	hoursCount                   = *JWT_SECRET_KEY_EXPIRE_HOUR
)

type Tokens struct {
	Access  string
	Refresh string
}

func GenerateNewTokens(id string, role string) (*Tokens, error) {
	accessToken, err := generateNewAccessToken(id, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateNewAccessToken(id string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix(),
	})
	return token.SignedString(privateKey)
}

func generateNewRefreshToken() (string, error) {
	hash := sha256.New()
	refresh := *JWT_SECRET_KEY + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime
	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
