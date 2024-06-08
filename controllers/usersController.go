package controllers

import (
	requestDto "my-go/dto/request"
	responseDto "my-go/dto/response"
	"my-go/initializers"
	"my-go/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	body := new(requestDto.LoginRequestDto)

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Failed to read body.",
			Data:    gin.H{},
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password.",
			Data:    gin.H{},
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Failed to create user.",
			Data:    gin.H{},
		})
	}

	c.JSON(http.StatusOK, responseDto.CustomResponseDto{
		Status:  http.StatusOK,
		Message: "success",
		Data:    gin.H{},
	})
}

func Login(c *gin.Context) {
	body := requestDto.LoginRequestDto{}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Failed to read body.",
			Data:    gin.H{},
		})
		return
	}

	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password.",
			Data:    gin.H{},
		})
		return
	}

	// Compare sent in password with saved users password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password.",
			Data:    gin.H{},
		})
		return
	}

	// Generate a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"role": "admin",
	})

	// Sign and get using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, responseDto.CustomResponseDto{
			Status:  http.StatusBadRequest,
			Message: "Failed to create token.",
			Data:    gin.H{},
		})
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, responseDto.CustomResponseDto{
		Status:  http.StatusOK,
		Message: "success",
		Data:    gin.H{},
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, responseDto.CustomResponseDto{
		Status:  http.StatusOK,
		Message: "success",
		Data:    user,
	})
}
