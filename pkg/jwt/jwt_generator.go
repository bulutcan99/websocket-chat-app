package jwt

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
)

type Tokens struct {
	Access  string
	Refresh string
}

func GenerateNewTokens(role string) (*Tokens, error) {
	accessToken, err := generateNewAccessToken(role)
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

func generateNewAccessToken(role string) (string, error) {
	secret := *JWT_SECRET_KEY
	minutesCount := *JWT_SECRET_KEY_EXPIRE_MINUTE

	claims := jwt.MapClaims{}
	claims["role"] = role
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
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
