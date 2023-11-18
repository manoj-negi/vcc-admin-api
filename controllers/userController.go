package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	//"strconv"
	//"log"
	"fmt"
	"net/http"
	"text/template"
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/models"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleName string `json:"rolename"`
	Mobile   string   `json:"mobile"`
}

func HandlerUserListing(c *fiber.Ctx) error {
	db, err := initializers.ConnectToDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error connecting to the database",
		})
	}

	query := `
		SELECT users.id, users.username, users.email,users.mobile,user_roles.name as role_name
		FROM users
		LEFT JOIN user_roles ON user_roles.id = users.role_id` 

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error executing the query",
		})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Mobile,
			&user.RoleName,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error scanning rows",
			})
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": users,
	})
}




func HandlerCreateUser(c *fiber.Ctx) error {
    db, _ := initializers.ConnectToDB()

    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Unable to hash password",
        })
    }

    _, err = db.Exec(context.Background(),
        "INSERT INTO users(username, role_id, api_key, client_id, country_code, email, password, validation_token, mobile, referral_code, product_id, total_invitees, successful_referral, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
        user.Username, user.RoleID, user.ApiKey, user.ClientID, user.CountryCode, user.Email, string(hash), user.ValidationToken, user.Mobile, user.ReferralCode, user.ProductID, user.TotalInvitees, user.SuccessfulReferral, user.IsActive)
    if err != nil {
        return err
    }

    err = db.QueryRow(context.Background(), "SELECT lastval()").Scan(&user.ID)
    if err != nil {
        return err
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": user,
    })
}

func Handleruser(c *fiber.Ctx) error{
	tmpl, err := template.ParseFiles("usercreate.html")
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
}

// func HandlerCreateUser(c *fiber.Ctx) error {
//     db, _ := initializers.ConnectToDB()

//      var user models.User
//     // if err := c.BodyParser(&user); err != nil {
//     //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//     //         "error": "Invalid request body",
//     //     })
//     // }

//     // For www-urlencoded data, use c.FormValue to get individual form values
//     Username := c.FormValue("username")
// 	fmt.Println("---userrrr---",Username)
// 	RoleID, _ := strconv.Atoi(c.FormValue("role_id"))
// 	fmt.Println("---role---id",RoleID)
// 	ApiKey := c.FormValue("api_key")
// 	ClientID := c.FormValue("client_id")
// 	CountryCode, _ := strconv.Atoi(c.FormValue("country_id"))
// 	Email := c.FormValue("email")
// 	Password := c.FormValue("password")
// 	ValidationToken := c.FormValue("validation_token")
// 	Mobile := c.FormValue("mobile")
// 	ReferralCode := c.FormValue("referral_code")
// 	ProductID, _ := strconv.Atoi(c.FormValue("product_id"))
// 	// TotalInvitees := c.FormValue("total_invitees")
// 	// SuccessfulReferral := c.FormValue("successful_referral")
// 	// IsActive := c.FormValue("is_active")

//     hash, err := bcrypt.GenerateFromPassword([]byte(Password), 10)
//     if err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "error": "Unable to hash password",
//         })
//     }

//     _, err = db.Exec(context.Background(),
//         "INSERT INTO users(username, role_id, api_key,client_id,country_id,email,password,validation_token,mobile,referral_code,product_id,total_invitees,successful_referral,is_active) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)",
//         Username,RoleID,ApiKey,ClientID,CountryCode,Email,string(hash),ValidationToken,Mobile,ReferralCode,ProductID,0,0,1)

//     if err != nil {
// 		fmt.Println("---errr-rrrr-------",err)
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "error": "Error inserting data into the database",
//         })
//     }

//     err = db.QueryRow(context.Background(), "SELECT lastval()").Scan(&user.ID)
//     if err != nil {
//         log.Fatal(err)
//     }

//     return c.Status(fiber.StatusCreated).JSON(fiber.Map{
//         "message": user,
//     })
// }

func HandlerGetUsersPaginatedPost(c *fiber.Ctx) error {
	// Define a struct to represent the request body
	var request struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}

	// Parse the request body into the defined struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate the request parameters
	if request.Page < 1 || request.PerPage < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page or per_page parameter",
		})
	}

	// Calculate the offset based on the page and perPage values
	offset := (request.Page - 1) * request.PerPage

	// Connect to the database
	db, err := initializers.ConnectToDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error connecting to the database",
		})
	}

	// Query to retrieve paginated users
	query := `
		SELECT id, username, email
		FROM users
		ORDER BY id
		OFFSET $1 LIMIT $2
	`

	rows, err := db.Query(context.Background(), query, offset, request.PerPage)
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error executing the query",
		})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
		); err != nil {
			fmt.Println("Error scanning rows:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error scanning rows",
			})
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":       request.Page,
		"per_page":   request.PerPage,
		"total_rows": len(users), // This should be the total number of rows in your table
		"users":      users,
	})
}


func HandlerGetAllUser(c *fiber.Ctx) error {
    db, _ := initializers.ConnectToDB()

    query := `SELECT
    id,
    username,
    role_id,
    api_key,
    client_id,
    country_code,
    email,
    validation_token,
    mobile,
    referral_code,
    product_id,
    total_invitees,
    successful_referral,
    is_active
FROM users;
`
    row, err := db.Query(context.Background(), query)

    if err != nil {
		fmt.Println("---------error",err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "error in database in users",
        })
    }

    var users []models.User
    defer row.Close()

    for row.Next() {
        var user models.User
        if err = row.Scan(
            &user.ID,
            &user.Username,
            &user.RoleID,
            &user.ApiKey,
            &user.ClientID,
            &user.CountryCode,
            &user.Email,
            &user.ValidationToken,
            &user.Mobile,
            &user.ReferralCode,
            &user.ProductID,
            &user.TotalInvitees,
            &user.SuccessfulReferral,
            &user.IsActive,
        ); err != nil {
            fmt.Println("errorrrr", err)
            return c.Status(http.StatusBadRequest).JSON(fiber.Map{
                "error": "error from getting the data from the database in users",
            })
        }
        users = append(users, user)
    }

    fmt.Println("ROwwwws", users)

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message":users,
    })
}

func generateRandomCode(length int) (string, error) {
	// Calculate the number of bytes needed for the random code
	numBytes := (length * 6) / 8

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomCode := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim the random code to the desired length
	randomCode = randomCode[:length]

	return randomCode, nil
}

// func UpdateOneUser(c *fiber.Ctx) error {
// 	db, _ := initializers.ConnectToDB()
// 	id := c.Params("id")

// 	var body models.User

// 	if err := c.BodyParser(&body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to read body",
// 		})
// 	}

// 	query := `
// 	UPDATE users
// 	SET first_name = $1, last_name = $2, email = $3, password = $4, age = $5, phone_no = $6
// 	WHERE id = $7
// 	 `

// 	result, err := db.Exec(context.Background(), query,
// 		body.First_name,
// 		body.Last_name,
// 		body.Email,
// 		body.Password,
// 		body.Age,
// 		body.Phone_no,
// 		id,
// 	)

// 	if err != nil {
// 		// Handl
// 		return err
// 	}
// 	rowsAffected := result.RowsAffected()

// 	if rowsAffected == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "User not found",
// 		})
// 	}

// 	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
// 		"message": "User Updated",
// 	})

// 	return nil
// }

func HandleUpdateUser(c *fiber.Ctx) error {
    db, _ := initializers.ConnectToDB()
    id := c.Params("id")

    var body models.User

    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to read body",
        })
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Unable to hash password",
        })
    }

    tx, err := db.Begin(context.Background())
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())

    query := `
        UPDATE users
        SET username = $1, role_id = $2, api_key = $3, client_id = $4, country_code = $5,
        email = $6, password = $7, validation_token = $8, mobile = $9, referral_code = $10,
        product_id = $11, total_invitees = $12, successful_referral = $13, is_active = $14
        WHERE id = $15
    `

    _, err = tx.Exec(context.Background(), query,
        body.Username,
        body.RoleID,
        body.ApiKey,
        body.ClientID,
        body.CountryCode,
        body.Email,
        string(hash), // Use the hashed password in the update
        body.ValidationToken,
        body.Mobile,
        body.ReferralCode,
        body.ProductID,
        body.TotalInvitees,
        body.SuccessfulReferral,
        body.IsActive,
        id,
    )
    if err != nil {
        return err
    }

    err = tx.Commit(context.Background())
    if err != nil {
        return err
    }

    c.Status(fiber.StatusAccepted).JSON(fiber.Map{
        "message": "User Updated",
    })

    return nil
}


func HandleDeleteUser(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	id := c.Params("id")

	query := `Delete from users where id = $1`

	result, err := db.Exec(context.Background(), query, id)

	if err != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "unable to execute the query",
		})
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Deleted",
	})

	return nil
}
