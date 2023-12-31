package controller

import (
	"github.com/bulutcan99/go-websocket/internal/model"
	"github.com/bulutcan99/go-websocket/internal/platform/cache"
	"github.com/bulutcan99/go-websocket/internal/repository"
	"github.com/bulutcan99/go-websocket/pkg/helper"
	"github.com/bulutcan99/go-websocket/pkg/token"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var (
	userRole = "user"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

	err := helper.EmailValidator(signUp.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Email is not valid",
		})
	}

	err = helper.PasswordValidator(signUp.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Password is not valid",
		})
	}

	err = helper.RoleValidator(userRole)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Role must be admin or user",
		})
	}

	user := model.User{
		UUID:         uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		BlockedAt:    time.Time{},
		UserName:     signUp.Name,
		UserSurName:  signUp.Surname,
		Nickname:     signUp.Nickname,
		Email:        signUp.Email,
		Passwordhash: utility.GeneratePassword(signUp.Password),
		Status:       1,
		UserRole:     userRole,
	}

	if errCreate := ac.repo.CreateUser(user); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to create user",
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

		isComparedUserPass := utility.ComparePasswords(getUser.Passwordhash, signIn.Password)
		if !isComparedUserPass {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Given password is not correct",
			})
		}

		role := getUser.UserRole
		userId := strconv.Itoa(int(getUser.Id))
		userIp := c.IP()
		tokens, err := token.GenerateNewTokens(userId, role, getUser.Email, userIp)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "There is a problem while trying to generate new tokens",
			})
		}

		err = ac.redisCache.SetAllUserData(signIn.Email, userId, userIp, getUser, tokens)
		err = ac.redisCache.SetUserIp(userId, userIp)
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   "Logged In Successfully!",
			"tokens": fiber.Map{
				"userid":  userId,
				"role":    role,
				"email":   signIn.Email,
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
		ipCache, err := ac.redisCache.GetUserIp(userDataWithCache.UserID)
		if err != nil && err == redis.Nil {
			tokens, err = token.GenerateNewTokens(userDataWithCache.UserID, userDataWithCache.UserRole, userDataWithCache.UserEmail, ipCache)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"msg":   "There is a problem while trying to generate new tokens",
				})
			}
		}

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
	err, tokenMetaData := ac.TokenProtection(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "User not authorized!",
		})
	}

	if tokenMetaData == nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "Token metadata is nil!",
		})
	}

	err = ac.redisCache.DeleteAllUserData(tokenMetaData.Email, tokenMetaData.Id)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to delete user data!",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged Out Successfully!",
	})
}

func (ac *AuthController) TokenProtection(c *fiber.Ctx) (error, *token.TokenMetaData) {
	now := time.Now().Unix()

	tokenMetaData, err := token.ExtractTokenMetaData(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to extract token metadata",
		}), nil
	}

	expires := tokenMetaData.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized! Check the expiration time of your token!",
		}), nil
	}

	userIpCache, err := ac.redisCache.GetUserIp(tokenMetaData.Id)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "There is an error while trying to get user token data",
		}), nil
	}

	if userIpCache != tokenMetaData.IP {
		return c.JSON(fiber.Map{
			"error": true,
			"msg":   "User ip is not same with token ip!",
		}), nil
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User is authorized!",
	}), tokenMetaData
}
