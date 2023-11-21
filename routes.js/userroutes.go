package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/controllers"
)

func SetupUserRoutes(app *fiber.App) {

	app.Post("/api/v1/login", controllers.HandlerAdminLogin)

	app.Get("/api/v1/country", controllers.HandlerGetAllCountry)

	app.Get("/api/v1/role", controllers.HandlerGetAllRole)

	app.Get("/api/v1/product", controllers.HandlerGetAllProduct)
	
	app.Get("/api/v1/user", controllers.HandlerGetAllUser)
	app.Post("/api/v1/user", controllers.HandlerCreateUser)
	app.Get("/api/v1/finduser/:id", controllers.HandlerGetOneUser)
	app.Get("/api/v1/userlisting", controllers.HandlerUserListing)
	app.Get("/api/v1/users", controllers.HandlerUserPagination)
	app.Put("/api/v1/user/:id", controllers.HandleUpdateUser)
	app.Delete("/api/v1/user/:id", controllers.HandleDeleteUser)
}
