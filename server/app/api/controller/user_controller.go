package controller

import (
	"github.com/bulutcan99/go-websocket/internal/db/cache"
	"github.com/bulutcan99/go-websocket/internal/db/repository"
	"github.com/bulutcan99/go-websocket/internal/model"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserInterface interface {
	GetUserSelf(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
	UserSignOut(c *fiber.Ctx) error
}

type UserController struct {
	repo       *repository.UserRepo
	redisCache *db_cache.RedisCache
	ac         *AuthController
}

func NewUserController(userRepo *repository.UserRepo, redisC *db_cache.RedisCache, authCont *AuthController) *UserController {
	return &UserController{
		repo:       userRepo,
		redisCache: redisC,
		ac:         authCont,
	}
}

func (uc *UserController) GetUserSelfInfo(c *fiber.Ctx) error {
	err, tokenMetaData := uc.ac.TokenProtection(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "User not authorized!",
		})
	}

	userDataWithCache, err := uc.redisCache.GetUserDataById(tokenMetaData.ID)
	if err == redis.Nil {
		id, err := uuid.Parse(tokenMetaData.ID)
		user, err := uc.repo.GetUserSelf(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"error": false,
			"user":  user,
		})
	} else {

		zap.S().Info("User data already in redis cache!")
		return c.JSON(fiber.Map{
			"error": false,
			"user":  userDataWithCache,
		})
	}
}

func (uc *UserController) GetAnotherUserInfo(c *fiber.Ctx) error {
	err, _ := uc.ac.TokenProtection(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "User not authorized!",
		})
	}

	email := c.Params("email")

	user, err := uc.repo.GetShowAnotherUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to get user info!",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"user":  user,
	})
}

func (uc *UserController) UpdatePasswordHandler(c *fiber.Ctx) error {
	err, tokenMetaData := uc.ac.TokenProtection(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "User not authorized!",
		})
	}

	newPass := &model.PasswordUpdate{}
	if err := c.BodyParser(newPass); err != nil {
		return custom_error.ParseError()
	}

	id, errUuid := uuid.Parse(tokenMetaData.ID)
	if errUuid != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to parse uuid",
		})
	}

	newPassHash := utility.GeneratePassword(newPass.NewPassword)
	updateErr := uc.repo.UpdatePassword(id, newPassHash)
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to update password",
		})
	}

	errUpdateRedis := uc.redisCache.UpdateUserPasswordHash(tokenMetaData.ID, newPassHash)
	if errUpdateRedis != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to update redis cache",
		})
	}

	return c.JSON(fiber.Map{
		"error":        false,
		"New Password": newPass,
	})
}
