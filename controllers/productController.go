package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	//"github.com/keshav/fiber/models"
)

type Product struct{
	ID   int   `json:"id"`
	Name string `json:"name"`
}


func HandlerGetAllProduct(c *fiber.Ctx) error {
	db, err := initializers.ConnectToDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error connecting to the database",
		})
	}

	query := `select id,name from products` 

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error executing the query",
		})
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error scanning rows",
			})
		}
		products = append(products, product)
	}

	return c.Status(fiber.StatusOK).JSON(products)
}


// func HandlerGetAllProduct(c *fiber.Ctx) error {
// 	db, err := initializers.ConnectToDB()
// 	if err != nil {
// 		fmt.Println("Error connecting to the database:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Error connecting to the database",
// 			StatusCode: fiber.StatusInternalServerError,
// 		})
// 	}

// 	query := `select id, name from products`

// 	rows, err := db.Query(context.Background(), query)
// 	if err != nil {
// 		fmt.Println("Error executing the query:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Error executing the query",
// 			StatusCode: fiber.StatusInternalServerError,
// 		})
// 	}
// 	defer rows.Close()

// 	var products []Product
// 	for rows.Next() {
// 		var product Product
// 		if err := rows.Scan(
// 			&product.ID,
// 			&product.Name,
// 		); err != nil {
// 			fmt.Println("Error scanning rows:", err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
// 				Status:     false,
// 				Message:    "Error scanning rows",
// 				StatusCode: fiber.StatusInternalServerError,
// 			})
// 		}
// 		products = append(products, product)
// 	}

// 	if len(products) > 0 {
// 		return c.Status(fiber.StatusOK).JSON(JsonResponse{
// 			Status:     true,
// 			Message:    "Data found",
// 			Data:       products,
// 			StatusCode: fiber.StatusOK,
// 		})
// 	} else {
// 		return c.Status(fiber.StatusNotFound).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Data not found",
// 			Data:       nil,
// 			StatusCode: fiber.StatusNotFound,
// 		})
// 	}
// }
