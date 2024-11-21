package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xelathan/mini_search_engine/db"
	"github.com/xelathan/mini_search_engine/utils"
	"github.com/xelathan/mini_search_engine/views"
)

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type LoginPayload struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := LoginPayload{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("Internal Server Error")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("Invalid Credentials")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(500)
		return c.SendString("Internal Server Error")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}

	token, err := jwt.ParseWithClaims(cookie, &utils.AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", 302)
	}

	_, ok := token.Claims.(*utils.AuthClaim)
	if !ok || !token.Valid {
		return c.Redirect("/login", 302)
	}

	return c.Next()
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSetting{}
	if err := settings.Get(); err != nil {
		return c.SendStatus(500)
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNewUrls))
}

func DashboardPostHandler(c *fiber.Ctx) error {
	settings := SettingsForm{}
	if err := c.BodyParser(&settings); err != nil {
		return c.SendStatus(500)
	}

	addNew := false
	if settings.AddNewUrls == "on" {
		addNew = true
	}

	searchOn := false
	if settings.SearchOn == "on" {
		searchOn = true
	}

	newSettings := &db.SearchSetting{}
	newSettings.Amount = settings.Amount
	newSettings.SearchOn = searchOn
	newSettings.AddNewUrls = addNew

	if err := newSettings.Update(); err != nil {
		return c.SendStatus(500)
	}

	c.Append("HX-Refresh", "true")

	return c.SendStatus(200)
}
