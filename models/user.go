package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := db.Create(user)
	if result != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, user *User) (string, error) {
	// Get user from email
	var userDB User
	result := db.Where("email = ?", user.Email).First(&userDB)

	if result.Error != nil {
		return "", result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	secretKey := "1234"

	// Create token
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["email"] = userDB.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
