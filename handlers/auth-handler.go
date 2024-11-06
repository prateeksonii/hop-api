package handlers

import (
	"crypto/subtle"
	"drop/db"
	"drop/models"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

type SignUpDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c echo.Context) error {
	dto := &SignUpDto{}
	if err := c.Bind(&dto); err != nil {
		return err
	}
	user := &models.User{}
	db.DB.Limit(1).Find(&user, "email = ?", dto.Email)

	if user.ID > 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "user exists with the same email",
		})
	}

	hash := argon2.IDKey([]byte(dto.Password), []byte("restinpeace"), 1, 64*1024, 4, 32)
	user.Name = dto.Name
	user.Email = dto.Email
	user.Password = base64.RawStdEncoding.EncodeToString(hash)

	res := db.DB.Create(&user)
	if err := res.Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &user)
}

type SignInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateTokens(userId uint) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"iss": "drop",
		"sub": strconv.Itoa(int(userId)),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte("restinpeace"))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}

	authSession := &models.AuthSession{}

	db.DB.Limit(1).Find(&authSession, "user_id = ?", userId)
	if authSession.ID > 0 {
		authSession.RefreshToken = refreshToken.String()
		authSession.ExpirationUtc = time.Now().UTC().Add(time.Hour * 24 * 365 * 10)
	} else {
		authSession = &models.AuthSession{
			UserID:        userId,
			RefreshToken:  refreshToken.String(),
			ExpirationUtc: time.Now().UTC().Add(time.Hour * 24 * 365 * 10),
		}
	}

	res := db.DB.Save(&authSession)
	if res.Error != nil {
		return "", "", res.Error
	}

	return tokenString, authSession.RefreshToken, nil
}

func SignIn(c echo.Context) error {
	dto := &SignInDto{}
	if err := c.Bind(&dto); err != nil {
		return err
	}

	user := &models.User{}
	db.DB.Limit(1).Find(&user, "email = ?", dto.Email)
	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	hash := argon2.IDKey([]byte(dto.Password), []byte("restinpeace"), 1, 64*1024, 4, 32)
	originalHash, err := base64.RawStdEncoding.DecodeString(user.Password)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(hash, originalHash) != 1 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	tokenString, refreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:   "hop_session",
		Value:  refreshToken,
		Path:   "/",
		Secure: false,
	})

	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	return c.JSON(http.StatusOK, &user)
}

func RefreshAuth(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	log.Println(authHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("invalid token")
	}

	tokenString := strings.Split(authHeader, "Bearer ")[1]

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("restinpeace"), nil
	})
	if err != nil {
		if err.Error() != "token has invalid claims: token is expired" {
			return err
		}
	}

	subject := ""
	for key, val := range claims {
		if key == "sub" {
			subject = val.(string)
		}
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		return err
	}

	session, err := c.Cookie("hop_session")
	if err != nil {
		return err
	}

	if err := session.Valid(); err != nil {
		return err
	}

	token, refreshToken, err := GenerateTokens(uint(userId))
	if err != nil {
		return err
	}

	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	c.SetCookie(&http.Cookie{
		Name:   "hop_session",
		Value:  refreshToken,
		Path:   "/",
		Secure: false,
	})

	return c.NoContent(http.StatusOK)
}

func GetAuthenticatedUser(c echo.Context) error {
	userId := GetUserId(c)
	user := &models.User{}

	db.DB.Find(&user, userId)
	return c.JSON(http.StatusOK, user)
}

func GetUserId(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	subject := ""

	for key, val := range claims {
		if key == "sub" {
			subject = val.(string)
		}
	}

	sub, err := strconv.Atoi(subject)
	if err != nil {
		return 0
	}
	return uint(sub)
}
