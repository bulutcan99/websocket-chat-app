package db_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bulutcan99/go-websocket/internal/model"
	config_redis "github.com/bulutcan99/go-websocket/pkg/config/redis"
	"github.com/bulutcan99/go-websocket/pkg/token"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type UserCache struct {
	UserID           string `json:"id"`
	UserUUID         string `json:"uuid"`
	UserName         string `json:"name"`
	UserSurname      string `json:"surname"`
	UserNickname     string `json:"nickname"`
	UserEmail        string `json:"email"`
	UserRole         string `json:"role"`
	UserCreatedAt    string `json:"created_at"`
	UserUpdatedAt    string `json:"updated_at"`
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

func NewRedisCache(redis *config_redis.Redis) *RedisCache {
	return &RedisCache{
		client:  redis.Client,
		context: redis.Context,
	}
}

func (rc *RedisCache) GetUserIp(id string) (string, error) {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:ip", id)
	userIp, err := rc.client.Get(rc.context, cacheUserIdKey).Result()
	if err != nil {
		return "", err
	}

	return userIp, nil
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

	return rc.GetUserDataById(*id)
}

func (rc *RedisCache) getUserId(email string) (*string, error) {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s:id", email)
	userID, err := rc.client.Get(rc.context, cacheUserEmailKey).Result()
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (rc *RedisCache) GetUserDataById(id string) (*UserCache, error) {
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

func (rc *RedisCache) SetAllUserData(email, id, ip string, user *model.User, tokens *token.Tokens) error {
	err := rc.setUserId(email, id)
	if err != nil {
		return err
	}

	cacheData := rc.SetDbUserToCacheUser(user)
	err = rc.SetUserData(id, cacheData)
	if err != nil {
		return err
	}

	err = rc.SetUserToken(id, tokens)
	if err != nil {
		return err
	}

	err = rc.SetUserIp(id, ip)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetUserIp(id string, ip string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:ip", id)
	err := rc.client.Set(rc.context, cacheUserIdKey, ip, 24*time.Hour).Err()
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

func (rc *RedisCache) SetDbUserToCacheUser(user *model.User) *UserCache {
	return &UserCache{
		UserID:           strconv.Itoa(int(user.Id)),
		UserUUID:         user.UUID.String(),
		UserName:         user.UserName,
		UserSurname:      user.UserSurName,
		UserNickname:     user.Nickname,
		UserPasswordHash: user.Passwordhash,
		UserEmail:        user.Email,
		UserRole:         user.UserRole,
		UserStatus:       strconv.Itoa(user.Status),
		UserCreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UserUpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
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

	err = rc.DeleteUserToken(id)
	if err != nil {
		return err
	}

	err = rc.deleteUserIp(id)

	return nil
}

func (rc *RedisCache) deleteUserIp(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:ip", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
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

func (rc *RedisCache) DeleteUserToken(id string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:token", id)
	err := rc.client.Del(rc.context, cacheUserIdKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) UpdateUserPasswordHash(id string, newPassHash string) error {
	cacheUserIdKey := fmt.Sprintf("user:id:%s:data", id)
	userData, err := rc.client.Get(rc.context, cacheUserIdKey).Result()
	if err != nil {
		return err
	}

	var user UserCache
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return err
	}

	user.UserPasswordHash = newPassHash
	user.UserUpdatedAt = time.Now().String()
	updatedUserData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rc.client.Set(rc.context, cacheUserIdKey, updatedUserData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}
