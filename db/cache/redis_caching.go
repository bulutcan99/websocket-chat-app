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

func (rc *RedisCache) GetUserDataByEmail(email string) (*UserCache, error) {
	id, err := rc.getUserId(email)
	if err != nil {
		return nil, err
	}

	return rc.getUserData(*id)
}

func (rc *RedisCache) getUserId(email string) (*string, error) {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	userID, err := rc.client.Get(rc.context, cacheUserEmailKey).Result()
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (rc *RedisCache) getUserData(id string) (*UserCache, error) {
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

func (rc *RedisCache) SetAllUserData(email string, id string, user *model.User, tokens *token.Tokens) error {
	err := rc.setUserId(email, id)
	if err != nil {
		return err
	}

	cacheData := rc.dbUserToCacheUser(user)
	err = rc.setUserData(id, cacheData)
	if err != nil {
		return err
	}

	err = rc.SetUserToken(id, tokens)
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

func (rc *RedisCache) setUserId(email string, id string) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	err := rc.client.Set(rc.context, cacheUserEmailKey, id, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) setUserData(id string, user *UserCache) error {
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

func (rc *RedisCache) dbUserToCacheUser(user *model.User) *UserCache {
	return &UserCache{
		UserID:           user.ID.String(),
		UserEmail:        user.Email,
		UserRole:         user.UserRole,
		UserCreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UserUpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		UserNameSurname:  user.NameSurname,
		UserPasswordHash: user.PasswordHash,
		UserStatus:       strconv.Itoa(user.Status),
	}
}

func (rc *RedisCache) DeleteAllUserData(email string, id string) error {
	err := rc.deleteUserId(email)
	if err != nil {
		return err
	}

	err = rc.deleteUserData(id)
	if err != nil {
		return err
	}

	err = rc.deleteUserToken(id)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) deleteUserId(email string) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	err := rc.client.Del(rc.context, cacheUserEmailKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) deleteUserData(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:data", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) deleteUserToken(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:token", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
	if err != nil {
		return err
	}

	return nil
}
