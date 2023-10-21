package controller

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/db/cache"
	"github.com/bulutcan99/go-websocket/db/repository"
	"github.com/bulutcan99/go-websocket/model"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/helper"
	"github.com/bulutcan99/go-websocket/pkg/token"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type AuthInterface interface {
	UserRegister(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
	UserSignOut(c *fiber.Ctx) error
}

type AuthController struct {
	repo       *repository.AuthRepo
	redisCache *cache.RedisCache
}

func NewAuthController(authRepo *repository.AuthRepo, redisC *cache.RedisCache) *AuthController {
	return &AuthController{
		repo:       authRepo,
		redisCache: redisC,
	}
}

func (ac *AuthController) UserRegister(c *fiber.Ctx) error {
	signUp := &model.Register{}
	if err := c.BodyParser(signUp); err != nil {
		return custom_error.ParseError()
	}

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
	fmt.Println("3131")
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

	if errCreate := ac.repo.CreateUser(user); errCreate != nil {
		errMsg := fmt.Sprintf("There is an error while create: %v", errCreate)
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   errMsg,
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User created successfully!",
		"user":  user,
	})
}

func (ac *AuthController) UserLogin(c *fiber.Ctx) error {
	signIn := &model.SignIn{}
	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

	userDataWithCache, err := ac.redisCache.GetUserDataByEmail(signIn.Email)
	if err != nil {
		getUser, err := ac.repo.GetUserSignByEmail(signIn.Email)
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

		role := getUser.UserRole
		userId := getUser.ID.String()
		tokens, err := token.GenerateNewTokens(getUser.ID.String(), role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is a problem while trying to generate new tokens",
			})
		}

		err = ac.redisCache.SetAllUserData(signIn.Email, userId, getUser, tokens)
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged In Successfully!",
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		zap.S().Info("User data already in redis cache!")
		isComparedUserPass := utility.ComparePasswords(userDataWithCache.UserPasswordHash, signIn.Password)
		if !isComparedUserPass {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Given password is not correct",
			})
		}

		tokens, err := ac.redisCache.GetUserTokenData(userDataWithCache.UserID)
		if err != nil && err == redis.Nil {
			tokens, err = token.GenerateNewTokens(userDataWithCache.UserID, userDataWithCache.UserRole)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"msg":   "There is a problem while trying to generate new tokens",
				})
			}

			err = ac.redisCache.SetUserToken(userDataWithCache.UserID, tokens)
			if err != nil {
				zap.S().Error("Error while trying to set user token: ", err)
			}
		}

		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged In Successfully!",
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	}
}

func (ac *AuthController) UserLogOut(c *fiber.Ctx) error {
	tokenMetaData, err := token.ExtractTokenMetaData(c)
	err = ac.redisCache.DeleteAllUserData(tokenMetaData.Email, tokenMetaData.ID)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged Out Successfully!",
		})
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
