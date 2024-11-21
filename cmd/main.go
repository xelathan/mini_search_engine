package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
	"github.com/xelathan/mini_search_engine/db"
	"github.com/xelathan/mini_search_engine/routes"
	"github.com/xelathan/mini_search_engine/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app := fiber.New(fiber.Config{
		AppName:     "Mini Search Engine",
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())
	db.InitDB()
	routes.SetRoutes(app)
	utils.StartCronJobs()

	errChannel := make(chan error, 1)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal(err)
			errChannel <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-c:
		fmt.Println("Received signal:", sig)
		fmt.Println("Shutting down gracefully...")
		app.Shutdown()
	case err := <-errChannel:
		fmt.Println("Error:", err)
	}
}
