package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bulutcan99/go-websocket/model"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/token"
	"github.com/redis/go-redis/v9"
	"strconv"
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

type RedisCacheInterface interface {
	GetUserId(email string) (string, error)
	GetUserData(id string) (*UserCache, error)
	SetUserId(email string, id string) error
	SetUserData(id string, user UserCache) error
}

type RedisCache struct {
	client  *redis.Client
	context context.Context
}

func NewRedisCache(redis *config.Redis) *RedisCache {
	return &RedisCache{
		client:  redis.Client,
		context: redis.Context,
	}
}

func (rc *RedisCache) GetUserId(email string) (*string, error) {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	userID, err := rc.client.Get(rc.context, cacheUserEmailKey).Result()
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (rc *RedisCache) GetUserData(id string) (*UserCache, error) {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:data", id)
	userData, err := rc.client.Get(rc.context, cacheUserIdKey).Result()
	if err != nil {
		return nil, err
	}

	var user UserCache
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (rc *RedisCache) GetUserTokenData(id string) (*token.Tokens, error) {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:token", id)
	userData, err := rc.client.Get(rc.context, cacheUserIdKey).Result()
	if err != nil {
		return nil, err
	}

	var tokens token.Tokens
	if err := json.Unmarshal([]byte(userData), &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (rc *RedisCache) SetUserId(email string, id string) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	err := rc.client.Set(rc.context, cacheUserEmailKey, id, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetUserData(id string, user *UserCache) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:data", id)
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rc.client.Set(rc.context, cacheUserIdKey, userData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetUserToken(id string, tokens *token.Tokens) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:token", id)
	tokenData, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	err = rc.client.Set(rc.context, cacheUserIdKey, tokenData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) DeleteUserId(email string) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	err := rc.client.Del(rc.context, cacheUserEmailKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) DeleteUserData(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:data", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) DeleteUserToken(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:token", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) CacheUserData(user model.User) error {
	userCache := UserCache{
		UserID:           user.ID.String(),
		UserEmail:        user.Email,
		UserRole:         user.UserRole,
		UserCreatedAt:    user.CreatedAt.String(),
		UserUpdatedAt:    user.UpdatedAt.String(),
		UserNameSurname:  user.NameSurname,
		UserPasswordHash: user.PasswordHash,
		UserStatus:       strconv.Itoa(user.Status),
	}

	return rc.SetUserData(user.ID.String(), &userCache)
}
