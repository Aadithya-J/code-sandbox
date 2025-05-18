package handler

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Aadithya-J/code-sandbox/internal/db"
	"github.com/Aadithya-J/code-sandbox/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(jwtSecret) == 0 {
		log.Println("Warning: JWT_SECRET environment variable not set. Using default insecure secret.")
		jwtSecret = []byte("your-very-secret-key-please-change-in-prod")
	}
}

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func UserRegister(c *fiber.Ctx) error {
	var request UserRegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	email := strings.TrimSpace(request.Email)
	password := request.Password

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email and password are required")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to process registration")
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).SendString("Email already exists")
		}
		log.Printf("Failed to create user: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func UserLogin(c *fiber.Ctx) error {
	var request UserLoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	email := strings.TrimSpace(request.Email)
	password := request.Password

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email and password are required")
	}

	var user models.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
		}
		log.Printf("Database error during login for email %s: %v", email, result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString("Error logging in")
	}

	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing token for user %s: %v", email, err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to login")
	}

	return c.JSON(fiber.Map{"token": t})
}
