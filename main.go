package main

import (
	"context"
	"log"
	"os"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" 
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/routes.js"
)

func init() {
	initializers.LoadVariable()
	initializers.ConnectToDB()
}

func main() {
    db, _ := initializers.ConnectToDB()
    defer db.Close(context.Background())

    engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

    app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))
    routes.SetupUserRoutes(app)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("Server listening on port %s", port)
    log.Fatal(app.Listen(":" + port))
}
