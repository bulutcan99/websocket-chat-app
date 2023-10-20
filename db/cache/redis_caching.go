package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	UserID           string `json:"id"`
	UserEmail        string `json:"email"`
	UserRole         string `json:"role"`
	UserCreatedAt    string `json:"created_at"`
	UserUpdatedAt    string `json:"updated_at"`
	UserNameSurname  string `json:"name_surname"`
	UserPasswordHash string `json:"password_hash"`
	UserStatus       string `json:"status"`
}

type UserCacheToken struct {
	UserAccessToken  string `json:"access_token"`
	UserRefreshToken string `json:"refresh_token"`
}

type RedisCacheInterface interface {
	GetUserId(email string) (string, error)
	GetUserData(id string) (*UserCache, error)
	SetUserId(email string, id string) error
	SetUserData(id string, user UserCache) error
}

type RedisCache struct {
	RedisCache *redis.Client
}

func NewRedisCache(redisCache *redis.Client) *RedisCache {
	return &RedisCache{
		RedisCache: redisCache,
	}
}

func (rc *RedisCache) GetUserId(email string) (*string, error) {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s", email)
	userID, err := rc.RedisCache.Get(context.Background(), cacheUserEmailKey).Result()
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (rc *RedisCache) GetUserData(id string) (*UserCache, error) {
	cacheUserIdKey := fmt.Sprintf("user:id:%s", id)
	userData, err := rc.RedisCache.HGetAll(context.Background(), cacheUserIdKey).Result()
	if err != nil {
		return nil, err
	}

	user := &UserCache{
		UserID:           id,
		UserEmail:        userData["userEmail"],
		UserRole:         userData["userRole"],
		UserCreatedAt:    userData["userCreatedAt"],
		UserUpdatedAt:    userData["userUpdatedAt"],
		UserNameSurname:  userData["userNameSurname"],
		UserPasswordHash: userData["userPassword"],
		UserStatus:       userData["userStatus"],
	}

	return user, nil
}

func (rc *RedisCache) GetUserTokenData(id string) (*UserCacheToken, error) {
	cacheUserIdKey := fmt.Sprintf("user:id:token:%s", id)
	userData, err := rc.RedisCache.HGetAll(context.Background(), cacheUserIdKey).Result()
	if err != nil {
		return nil, err
	}

	user := &UserCacheToken{
		UserAccessToken:  userData["userAccessToken"],
		UserRefreshToken: userData["userRefreshToken"],
	}

	return user, nil
}

func (rc *RedisCache) SetUserId(email string, id string) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s", email)
	err := rc.RedisCache.Set(context.Background(), cacheUserEmailKey, id, time.Hour*12).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetUserData(id string, user UserCache) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s", id)
	userData := map[string]interface{}{
		"userEmail":       user.UserEmail,
		"userNameSurname": user.UserNameSurname,
		"userPassword":    user.UserPasswordHash,
		"userCreatedAt":   user.UserCreatedAt,
		"userUpdatedAt":   user.UserUpdatedAt,
		"userRole":        user.UserRole,
		"userStatus":      user.UserStatus,
	}

	err := rc.RedisCache.HMSet(context.Background(), cacheUserIdKey, userData).Err()
	if err != nil {
		return err
	}

	err = rc.RedisCache.Expire(context.Background(), cacheUserIdKey, 12*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetUserToken(id string, tokens UserCacheToken) error {
	cacheUserIdKey := fmt.Sprintf("user:id:token:%s", id)
	token := map[string]interface{}{
		"accessToken":  tokens.UserAccessToken,
		"refreshToken": tokens.UserRefreshToken,
	}

	err := rc.RedisCache.HMSet(context.Background(), cacheUserIdKey, token).Err()
	if err != nil {
		return err
	}

	err = rc.RedisCache.Expire(context.Background(), cacheUserIdKey, 12*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}
