package cache

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

type RedisCache struct {
	RedisCache *redis.Client
}

func NewRedisCache(redisCache *redis.Client) *RedisCache {
	return &RedisCache{
		RedisCache: redisCache,
	}
}

func (rc *RedisCache) GetUserDataFromCache(email string) (*UserCache, error) {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s", email)
	userID, err := rc.RedisCache.SMembers(context.Background(), cacheUserEmailKey).Result()
	if err == redis.Nil {
		return nil, err
	}

	cacheUserIdKey := fmt.Sprintf("user:id:%s", userID[0])
	userData, err := rc.RedisCache.HGetAll(context.Background(), cacheUserIdKey).Result()
	if err == redis.Nil {
		return nil, err
	}

	return &UserCache{
		UserID:           userID[0],
		UserEmail:        userData["userEmail"],
		UserRole:         userData["userRole"],
		UserCreatedAt:    userData["userCreatedAt"],
		UserUpdatedAt:    userData["userUpdatedAt"],
		UserNameSurname:  userData["userNameSurname"],
		UserPasswordHash: userData["userPassword"],
		UserStatus:       userData["userStatus"],
	}, nil
}

func (rc *RedisCache) SetUserCacheData(userCache UserCache) error {
	cacheUserEmailKey := fmt.Sprintf("user:email:%s", userCache.UserEmail)
	err := rc.RedisCache.SAdd(context.Background(), cacheUserEmailKey, userCache.UserID, time.Hour*12).Err()
	if err != nil {
		return err
	}

	cacheUserIdKey := fmt.Sprintf("user:id:%s", userCache.UserID)
	userData := map[string]interface{}{
		"userEmail":       userCache.UserEmail,
		"userNameSurname": userCache.UserNameSurname,
		"userPassword":    userCache.UserPasswordHash,
		"userCreatedAt":   userCache.UserCreatedAt,
		"userUpdatedAt":   userCache.UserUpdatedAt,
		"userRole":        userCache.UserRole,
		"userStatus":      userCache.UserStatus,
	}

	err = rc.RedisCache.HMSet(context.Background(), cacheUserIdKey, userData).Err()
	if err != nil {
		return err
	}

	err = rc.RedisCache.Expire(context.Background(), cacheUserIdKey, 12*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) SetRefreshTokenCookie(c *fiber.Ctx, refreshToken string) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
}
