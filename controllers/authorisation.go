package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"unicode"
)

func Register(c *fiber.Ctx) error {
	//Rendering page with registration form
	return c.Render("register", fiber.Map{
		"Title": "Registration page",
	})
}

func RegisterConfirm(c *fiber.Ctx) error {
	//Getting data from form
	var data struct {
		Username string
		Password string
	}
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	//Checking if Username is valid
	var validUchars, validUupper, validUlower bool
	validUchars = true
	for _, char := range data.Username {
		if unicode.IsLower(char) {
			validUlower = true
		}
		if unicode.IsUpper(char) {
			validUupper = true
		}
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			validUchars = false
		}
	}
	if validUchars == false || validUlower == false || validUupper == false || len(data.Username) > 40 || len(data.Username) < 5 {
		return c.Render("register", fiber.Map{
			"Error": "Invalid username",
		})
	}
	//Checking if Password is valid
	var validPchars, validPlower, validPupper, validPdigit bool
	validPchars = true
	for _, char := range data.Password {
		if unicode.IsUpper(char) {
			validPupper = true
		}
		if unicode.IsLower(char) {
			validPlower = true
		}
		if unicode.IsDigit(char) {
			validPdigit = true
		}
		if !unicode.IsDigit(char) && !unicode.IsLetter(char) {
			validPchars = false
		}
	}
	if validPlower == false || validPchars == false || validPupper == false || validPdigit == false || len(data.Password) > 20 || len(data.Password) < 8 {
		return c.Render("register", fiber.Map{
			"Error": "Invalid password",
		})
	}
	//Checking if Username is already taken
	result := database.DB.Where("username = ?", data.Username)
	if result.Error == gorm.ErrRecordNotFound {
		return c.Render("register", fiber.Map{
			"Error": "Username already taken",
		})
	}
	//Hashing password
	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: data.Username,
		Password: password,
	}
	//Creating user in database
	database.DB.Create(&user)
	return c.Redirect("/")
}

func Login(c *fiber.Ctx) error {
	//Rendering page with login form
	return c.Render("login", fiber.Map{
		"Title": "Login page",
	})
}

func LoginConfirm(c *fiber.Ctx) error {
	//Getting data from form
	var data struct {
		Username string
		Password string
	}
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	//Getting info about user with given username
	var user models.User
	result := database.DB.Where("username = ?", data.Username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return c.Render("login", fiber.Map{
			"Error": "Incorrect username or password",
		})
	}
	//Comparing hashed password and given password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.Render("login", fiber.Map{
			"Error": "incorrect username or password",
		})
	}
	//Creating new session
	database.CreateSession(c, &user)
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	//Deleting current session
	database.DeleteSession(c)
	return c.Redirect("/")
}
