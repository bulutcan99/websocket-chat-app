package controller

import (
	"fmt"
	db_cache "github.com/bulutcan99/go-websocket/db/cache"
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
	redisCache *db_cache.RedisCache
}

func NewAuthController(authRepo *repository.AuthRepo, redisC *db_cache.RedisCache) *AuthController {
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
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    tokens.Access,
			Expires:  time.Now().Add(time.Hour * 2),
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Lax",
		})

		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged In Successfully!",
			"tokens": fiber.Map{
				"userid":  userId,
				"role":    role,
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
		}

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    tokens.Access,
			Expires:  time.Now().Add(time.Hour * 12),
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Lax",
		})
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged In Successfully!",
			"tokens": fiber.Map{
				"userid":  userDataWithCache.UserID,
				"role":    userDataWithCache.UserRole,
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	}
}
func (ac *AuthController) UserLogOut(c *fiber.Ctx) error {
	now := time.Now().Unix()
	tokenMetaData, err := token.ExtractTokenMetaData(c)
	fmt.Println(tokenMetaData)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged Out Successfully!",
		})
	}

	err = ac.redisCache.DeleteAllUserData(tokenMetaData.Email, tokenMetaData.ID)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged Out Successfully!",
		})
	}

	expires := tokenMetaData.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized! Check the expiration time of your token",
		})
	}

	c.ClearCookie("access_token")

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged Out Successfully!",
	})
}
