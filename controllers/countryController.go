package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
)

type Country struct{
	ID   int   `json:"id"`
	Name string  `json:"name"`
}

type JsonResponse struct {
    Status     bool        `json:"status"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    StatusCode int         `json:"status_code"`
}


func HandlerGetAllCountry(c *fiber.Ctx) error {
	db, err := initializers.ConnectToDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error connecting to the database",
		})
	}

	query := `select code,name from countries` 

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error executing the query",
		})
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var country Country
		if err := rows.Scan(
			&country.ID,
			&country.Name,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error scanning rows",
			})
		}
		countries = append(countries, country)
	}

	return c.Status(fiber.StatusOK).JSON(countries)
}



// func HandlerGetAllCountry(c *fiber.Ctx) error {
//     db, err := initializers.ConnectToDB()
//     if err != nil {
//         fmt.Println("Error connecting to the database:", err)
//         return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error connecting to the database",
//             StatusCode: fiber.StatusInternalServerError,
//         })
//     }

//     query := `select code, name from countries`

//     rows, err := db.Query(context.Background(), query)
//     if err != nil {
//         fmt.Println("Error executing the query:", err)
//         return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error executing the query",
//             StatusCode: fiber.StatusInternalServerError,
//         })
//     }
//     defer rows.Close()

//     var countries []Country
//     for rows.Next() {
//         var country Country
//         if err := rows.Scan(
//             &country.ID,
//             &country.Name,
//         ); err != nil {
//             fmt.Println("Error scanning rows:", err)
//             return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//                 Status:     false,
//                 Message:    "Error scanning rows",
//                 StatusCode: fiber.StatusInternalServerError,
//             })
//         }
//         countries = append(countries, country)
//     }

//     if len(countries) > 0 {
//         return c.Status(fiber.StatusOK).JSON(JsonResponse{
//             Status:     true,
//             Message:    "Data found",
//             Data:       countries,
//             StatusCode: fiber.StatusOK,
//         })
//     } else {
//         return c.Status(fiber.StatusNotFound).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Data not found",
//             Data:       nil,
//             StatusCode: fiber.StatusNotFound,
//         })
//     }
// }

