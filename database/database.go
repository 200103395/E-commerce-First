package database

import (
	"Ecommerce/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var DB *gorm.DB
var secretKey = "NewSecretKey"

func Setup() {
	//Setting up mysql database with gorm
	connection, err := gorm.Open(mysql.Open("root:Ak200222!@tcp(localhost:3306)/test"), &gorm.Config{})
	//root = username, Ak200222! = password, test = db name
	if err != nil {
		panic("couldn't connect to database")
	}
	DB = connection
	//Migrating all needed models to database
	if err = connection.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.Item{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.Rating{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.Comment{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.RoleRequest{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.Cart{}); err != nil {
		panic(err)
	}
	if err = connection.AutoMigrate(&models.TokenRequest{}); err != nil {
		panic(err)
	}
}

func CreateSession(c *fiber.Ctx, user *models.User) error {
	//Creating new jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.UserID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	//Encoding token using secretKey
	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}
	//Creating cookie for new session
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func ReadSession(c *fiber.Ctx, user *models.User) string {
	//Getting and decoding jwt cookie
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusUnauthorized)
		return "Unauthorised"
	}
	//Getting information about user from token
	claims := token.Claims.(*jwt.StandardClaims)
	DB.Where("user_id = ?", claims.Issuer).First(user)

	return "Authorised"
}

func DeleteSession(c *fiber.Ctx) error {
	//Creating overlapping empty jwt token
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
