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
	"strconv"
	"time"
)

type AuthInterface interface {
	UserSignUp(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
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

func (ac *AuthController) UserSignUp(c *fiber.Ctx) error {
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
		return custom_error.DatabaseError()
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "User created successfully! User ID: " + user.ID.String(),
		"user":  user,
	})
}

func (ac *AuthController) UserSignIn(c *fiber.Ctx) error {
	// Basta redis ile emailden data id var mi diye kontrol et zaten id varsa
	// Diger datalar da olmus olacak ancak bunu rediste var mi yok mu diye iki farkli fonksiyona ayirarak burda cagir
	signIn := &model.SignIn{}
	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

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

	accessToken, err := token.GenerateNewTokens(getUser.ID.String(), role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to generate new tokens",
		})
	}

	userId := getUser.ID.String()
	errSaveToRedis := ac.redisCache.Set(userId, accessToken.Refresh, time.Minute*5).Err()
	if errSaveToRedis != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to save redis",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged In Successfully!",
		"tokens": fiber.Map{
			"access":  accessToken.Access,
			"refresh": accessToken.Refresh,
		},
	})
}

func createSuccessResponse(c *fiber.Ctx, tokens token.Tokens) error {
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged In Successfully!",
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

func (ac *AuthController) cacheUserData(user model.User) error {
	userCache := cache.UserCache{
		UserID:           user.ID.String(),
		UserEmail:        user.Email,
		UserRole:         user.UserRole,
		UserCreatedAt:    user.CreatedAt.String(),
		UserUpdatedAt:    user.UpdatedAt.String(),
		UserNameSurname:  user.NameSurname,
		UserPasswordHash: user.PasswordHash,
		UserStatus:       strconv.Itoa(user.Status),
	}
	return ac.redisCache.SetUserData(user.ID.String(), &userCache)
}

func (ac *AuthController) cacheUserTokens(user model.User, tokens token.Tokens) error {
	err := ac.redisCache.SetUserId(user.Email, user.ID.String())
	if err != nil {
		return err
	}

	return ac.redisCache.SetUserToken(user.ID.String(), tokens)
}

func (ac *AuthController) SignOut(c *fiber.Ctx) error {
	tokenMetaData, err := token.ExtractToken(c)
	deleted, err := ac.redisCache.Del(tokenMetaData.ID).Result()
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
