package controller

import (
	"context"
	"fmt"
	"github.com/bulutcan99/go-websocket/app/model"
	"github.com/bulutcan99/go-websocket/app/repository"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/helper"
	platform "github.com/bulutcan99/go-websocket/pkg/platform/cache"
	"github.com/bulutcan99/go-websocket/pkg/token"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type AuthInterface interface {
	UserSignUp(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
}

type AuthController struct {
	repo       *repository.AuthRepo
	redisCache *platform.RedisCache
}

func NewAuthController(authRepo *repository.AuthRepo, redisC *platform.RedisCache) *AuthController {
	return &AuthController{
		repo:       authRepo,
		redisCache: redisC,
	}
}

func (ac *AuthController) UserSignUp(c *fiber.Ctx) error {
	signUp := &model.Register{}
	if err := c.BodyParser(signUp); err != nil {
		return custom_error.ParseError()
	}

	signUpCtx := context.Background()
	err := helper.EmailValidator(signUp.Email)
	if err != nil {
		return fmt.Errorf("error while trying to set validation funcEmail, %w", err)
	}
	err = helper.PasswordValidator(signUp.Password)
	if err != nil {
		return fmt.Errorf("error while trying to set validation funcPass, %w", err)
	}
	err = helper.RoleValidator(signUp.UserRole)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Role must be admin or user",
		})
	}

	_, errVerify := utility.VerifyRole(signUp.UserRole)
	if errVerify != nil {
		return custom_error.ValidationError()
	}

	user := model.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		NameSurname:  signUp.NameSurname,
		Email:        signUp.Email,
		PasswordHash: utility.GeneratePassword(signUp.Password),
		Status:       1,
		UserRole:     signUp.UserRole,
	}

	if errCreate := ac.repo.CreateUser(signUpCtx, user); errCreate != nil {
		return custom_error.DatabaseError()
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User Created Successfully!",
		"user":  user.NameSurname,
	})
}

func (ac *AuthController) UserSignIn(c *fiber.Ctx) error {
	// Oncelikle mail ile rediste id var mi diye kontrol edilmeli, eger varsa id ile rediste user var mi diye kontrol edilmeli
	// Ardindan eger ikisi de varsa db cekmeye gerek yok, eger yoksa db cekilmeli ve redise kaydedilmeli
	signIn := &model.SignIn{}
	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

	getUser, err := ac.repo.GetUserSignByEmail(context.Background(), signIn.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Given email is not found",
		})
	}

	isComparedUserPass := utility.ComparePasswords(getUser.PasswordHash, signIn.Password)
	if !isComparedUserPass {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Given password is not correct",
		})
	}

	userCache := model.UserCache{
		UserID:           getUser.ID.String(),
		UserEmail:        getUser.Email,
		UserRole:         getUser.UserRole,
		UserCreatedAt:    getUser.CreatedAt.String(),
		UserUpdatedAt:    getUser.UpdatedAt.String(),
		UserNameSurname:  getUser.NameSurname,
		UserPasswordHash: getUser.PasswordHash,
		UserStatus:       strconv.Itoa(getUser.Status),
	}

	tokens, err := token.GenerateNewTokens(userCache.UserID, userCache.UserRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to generate new tokens",
		})
	}

	cacheUserEmailKey := fmt.Sprintf("user:email:%s", userCache.UserEmail)
	_, err = ac.redisCache.SMembers(context.Background(), cacheUserEmailKey).Result()
	if err == redis.Nil {
		fmt.Println("No data in for this user: ", userCache.UserEmail)
		errSaveToRedis := ac.redisCache.SAdd(context.Background(), cacheUserEmailKey, userCache.UserID, time.Hour*12).Err()
		if errSaveToRedis != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is a problem while trying to save email to redis",
			})
		}
	}

	cacheUserIdKey := fmt.Sprintf("user:id:%s", userCache.UserID)
	_, err = ac.redisCache.HGetAll(context.Background(), cacheUserIdKey).Result()
	if err == redis.Nil {
		fmt.Println("No data in for this user: ", userCache.UserID)
		userData := map[string]string{
			"userEmail":       userCache.UserEmail,
			"userNameSurname": userCache.UserNameSurname,
			"userPassword":    userCache.UserPasswordHash,
			"userCreatedAt":   userCache.UserCreatedAt,
			"userUpdatedAt":   userCache.UserUpdatedAt,
			"userRole":        userCache.UserRole,
			"userStatus":      userCache.UserStatus,
			"accessToken":     tokens.Access,
			"refreshToken":    tokens.Refresh,
		}

		err = ac.redisCache.HMSet(context.Background(), cacheUserIdKey, userData).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is a problem while trying to save user to redis",
			})
		}

		err = ac.redisCache.Expire(context.Background(), cacheUserIdKey, 12*time.Hour).Err()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("burda")
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokens.Refresh,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged In Successfully!",
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

func (ac *AuthController) SignOut(c *fiber.Ctx) error {
	tokenMetaData, err := token.ExtractToken(c)
	deleted, err := ac.redisCache.Del(context.Background(), tokenMetaData.ID).Result()
	if err != nil || deleted == 0 {
		return custom_error.ValidationError()
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "deleted cookie",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged Out Successfully!",
	})
}
