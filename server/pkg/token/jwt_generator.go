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
	JWT_SECRET_KEY             = &env.Env.JwtSecretKey
	JWT_SECRET_KEY_EXPIRE_TIME = &env.Env.JwtSecretKeyExpireTime
	JWT_SECRET_KEY_EXPIRE_HOUR = &env.Env.JwtRefreshKeyExpireHoursCount
)

type TokenMetaData struct {
	UUID    string
	Email   string
	Role    string
	Expires int64
}

type Tokens struct {
	Access  string
	Refresh string
}

func GenerateNewTokens(id string, role string, email string) (*Tokens, error) {
	accessToken, err := GenerateNewAccessToken(id, role, email)
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

func GenerateNewAccessToken(id string, role string, email string) (string, error) {
	timeCount := *JWT_SECRET_KEY_EXPIRE_TIME
	expUnix := time.Now().Add(time.Hour * time.Duration(timeCount)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"role":  role,
		"email": email,
		"exp":   expUnix,
	})

	return token.SignedString([]byte(*JWT_SECRET_KEY))
}

func generateNewRefreshToken() (string, error) {
	hash := sha256.New()
	refresh := *JWT_SECRET_KEY + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount := *JWT_SECRET_KEY_EXPIRE_HOUR
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime
	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
