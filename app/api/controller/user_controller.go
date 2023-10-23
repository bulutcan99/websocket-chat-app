package controller

import (
	db_cache "github.com/bulutcan99/go-websocket/db/cache"
	"github.com/bulutcan99/go-websocket/db/repository"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/token"
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
}

func NewUserController(userRepo *repository.UserRepo, redisC *db_cache.RedisCache) *UserController {
	return &UserController{
		repo:       userRepo,
		redisCache: redisC,
	}
}

func (uc *UserController) GetUserSelfHandler(c *fiber.Ctx) error {
	tokenMetaData, err := token.ExtractTokenMetaData(c)
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

		errSetRedis := uc.redisCache.SetUserData(userDataWithCache.UserID, userDataWithCache)
		if errSetRedis != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is an error while trying to set redis user data",
			})
		}

		return c.JSON(fiber.Map{
			"error": false,
			"user":  user,
		})
	} else {

		zap.S().Errorf("User data already in redis cache!")
		return c.JSON(fiber.Map{
			"error": false,
			"user":  userDataWithCache,
		})
	}
}

func (uc *UserController) GetAnotherUserHandler(c *fiber.Ctx) error {
	email := c.Params("email")

	user, err := uc.repo.GetShownedUserByEmail(email)
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
	tokenMetaData, err := token.ExtractTokenMetaData(c)
	var newPass *string
	if err := c.BodyParser(newPass); err != nil {
		return custom_error.ParseError()
	}
	userDataWithCache, err := uc.redisCache.GetUserDataById(tokenMetaData.ID)
	if err == redis.Nil {
		id, err := uuid.Parse(tokenMetaData.ID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "There is an error while trying to parse user id",
			})
		}

		updateErr := uc.repo.UpdatePassword(id, *newPass)
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is an error while trying to update password",
			})
		}

		return c.JSON(fiber.Map{
			"error":        false,
			"New Password": newPass,
		})
	} else {

		zap.S().Errorf("User data already in redis cache!")
		errSetRedis := uc.redisCache.SetUserData(userDataWithCache.UserID, userDataWithCache)
		if errSetRedis != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is an error while trying to set redis user data",
			})
		}
	}
	return nil
}
