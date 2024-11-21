package routes

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

type SettingsForm struct {
	Amount     int    `form:"amount"`
	SearchOn   string `form:"searchOn"`
	AddNewUrls string `form:"addNewUrls"`
}

func SetRoutes(app *fiber.App) {
	app.Get("/", AuthMiddleware, DashboardHandler)
	app.Post("/", AuthMiddleware, DashboardPostHandler)
	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Post("/search", HandleSearch)
	app.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))
}
