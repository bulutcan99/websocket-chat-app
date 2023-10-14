package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"
	"websocket_chat/pkg/env"
)

var (
	Env = env.ParseEnv()
)

type Tokens struct {
	Access  string
	Refresh string
}

func GenerateNewTokens(id string, role string) (*Tokens, error) {
	accessToken, err := generateNewAccessToken(id, role)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateNewAccessToken(id string, role string) (string, error) {
	secret := Env.JwtSecretKey
	minutesCount := Env.JwtSecretKeyExpireMinutesCount

	claims := jwt.MapClaims{}
	claims["id"] = id
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
	refresh := Env.JwtRefreshKey + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount := Env.JwtRefreshKeyExpireHoursCount
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime
	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
