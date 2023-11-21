package controllers

import (
	"context"
	"strconv"
	"fmt"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/auth"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/models"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID                 int      `json:"id"`
	Username           string   `json:"username"`
	RoleName           string   `json:"role_name"`
	ApiKey             string   `json:"api_key"`
	ClientID           string   `json:"client_id"`
	CountryName        string   `json:"country_name"`
	Email              string   `json:"email"`
	Password           string   `json:"password"`
	ValidationToken    string   `json:"validation_token"`
	Mobile             string   `json:"mobile"`
	ReferralCode       string   `json:"referral_code"`
	ProductName        string   `json:"product_name"`
	TotalInvitees      int      `json:"total_invitees"`
	SuccessfulReferral int      `json:"successful_referral"`
	IsActive           int      `json:"is_active"`
}


type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleName string `json:"rolename"`
	Mobile   string `json:"mobile"`
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
	SELECT users.id, users.username, users.email, users.mobile, user_roles.name as role_name
	FROM users
	LEFT JOIN user_roles ON user_roles.id = users.role_id
	ORDER BY created_at DESC
	LIMIT 10` 

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

// func HandlerUserListing(c *fiber.Ctx) error {
// 	db, err := initializers.ConnectToDB()
// 	if err != nil {
// 		fmt.Println("Error connecting to the database:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Error connecting to the database",
// 			StatusCode: fiber.StatusInternalServerError,
// 		})
// 	}

// 	query := `
// 	SELECT users.id, users.username, users.email, users.mobile, user_roles.name as role_name
// 	FROM users
// 	LEFT JOIN user_roles ON user_roles.id = users.role_id
// 	ORDER BY created_at DESC
// 	LIMIT 10`

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

// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(
// 			&user.ID,
// 			&user.Username,
// 			&user.Email,
// 			&user.Mobile,
// 			&user.RoleName,
// 		); err != nil {
// 			fmt.Println("Error scanning rows:", err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
// 				Status:     false,
// 				Message:    "Error scanning rows",
// 				StatusCode: fiber.StatusInternalServerError,
// 			})
// 		}
// 		users = append(users, user)
// 	}

// 	if len(users) > 0 {
// 		return c.Status(fiber.StatusOK).JSON(JsonResponse{
// 			Status:     true,
// 			Message:    "Data found",
// 			Data:       users,
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

func HandlerCreateUser(c *fiber.Ctx) error {
    db, _ := initializers.ConnectToDB()

    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

	hashedPassword := auth.Sha256Hash(user.Password)

    _, err := db.Exec(context.Background(),
        "INSERT INTO users(username, role_id, api_key, client_id, country_code, email, password, validation_token, mobile, referral_code, product_id, total_invitees, successful_referral, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
        user.Username, user.RoleID, user.ApiKey, user.ClientID, user.CountryCode, user.Email, string(hashedPassword), user.ValidationToken, user.Mobile, user.ReferralCode, user.ProductID, user.TotalInvitees, user.SuccessfulReferral, user.IsActive)
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

// func HandlerCreateUser(c *fiber.Ctx) error {
//     db, _ := initializers.ConnectToDB()

//     var user models.User
//     if err := c.BodyParser(&user); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Invalid request body",
//             StatusCode: fiber.StatusBadRequest,
//         })
//     }

//     hashedPassword := auth.Sha256Hash(user.Password)

//     _, err := db.Exec(context.Background(),
//         "INSERT INTO users(username, role_id, api_key, client_id, country_code, email, password, validation_token, mobile, referral_code, product_id, total_invitees, successful_referral, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
//         user.Username, user.RoleID, user.ApiKey, user.ClientID, user.CountryCode, user.Email, string(hashedPassword), user.ValidationToken, user.Mobile, user.ReferralCode, user.ProductID, user.TotalInvitees, user.SuccessfulReferral, user.IsActive)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error creating user",
//             StatusCode: fiber.StatusInternalServerError,
//         })
//     }

//     err = db.QueryRow(context.Background(), "SELECT lastval()").Scan(&user.ID)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error retrieving user ID",
//             StatusCode: fiber.StatusInternalServerError,
//         })
//     }

//     return c.Status(fiber.StatusCreated).JSON(JsonResponse{
//         Status:     true,
//         Message:    "User created successfully",
//         Data:       user,
//         StatusCode: fiber.StatusCreated,
//     })
// }

func HandlerUserPagination(c *fiber.Ctx) error {
    db, _ := initializers.ConnectToDB()

    page, err := strconv.Atoi(c.Query("page", "1"))
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(c.Query("limit", "10"))
    if err != nil || limit < 1 {
        limit = 10
    }
    
	offset := (page - 1) * limit

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
	FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

    rows, err := db.Query(context.Background(), query, limit, offset)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Error in database in users",
        })
    }
    defer rows.Close()

    var users []models.User

    for rows.Next() {
        var user models.User
        if err := rows.Scan(
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
            fmt.Println("Error scanning data:", err)
            return c.Status(http.StatusBadRequest).JSON(fiber.Map{
                "error": "Error getting data from the database in users",
            })
        }
        users = append(users, user)
    }

    return c.JSON(users)
}

// func HandlerUserPagination(c *fiber.Ctx) error {
//     db, _ := initializers.ConnectToDB()

//     page, err := strconv.Atoi(c.Query("page", "1"))
//     if err != nil || page < 1 {
//         page = 1
//     }

//     limit, err := strconv.Atoi(c.Query("limit", "10"))
//     if err != nil || limit < 1 {
//         limit = 10
//     }
    
// 	offset := (page - 1) * limit

//     query := `SELECT 
// 	id,
// 	username,
// 	role_id,
// 	api_key,
// 	client_id,
// 	country_code,
// 	email,
// 	validation_token,
// 	mobile,
// 	referral_code,
// 	product_id,
// 	total_invitees,
// 	successful_referral,
// 	is_active 
// 	FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

//     rows, err := db.Query(context.Background(), query, limit, offset)
//     if err != nil {
//         fmt.Println("Error executing query:", err)
//         return c.Status(http.StatusBadRequest).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error in database in users",
//             StatusCode: http.StatusBadRequest,
//         })
//     }
//     defer rows.Close()

//     var users []models.User

//     for rows.Next() {
//         var user models.User
//         if err := rows.Scan(
//             &user.ID,
//             &user.Username,
//             &user.RoleID,
//             &user.ApiKey,
//             &user.ClientID,
//             &user.CountryCode,
//             &user.Email,
//             &user.ValidationToken,
//             &user.Mobile,
//             &user.ReferralCode,
//             &user.ProductID,
//             &user.TotalInvitees,
//             &user.SuccessfulReferral,
//             &user.IsActive,
//         ); err != nil {
//             fmt.Println("Error scanning data:", err)
//             return c.Status(http.StatusBadRequest).JSON(JsonResponse{
//                 Status:     false,
//                 Message:    "Error getting data from the database in users",
//                 StatusCode: http.StatusBadRequest,
//             })
//         }
//         users = append(users, user)
//     }

//     return c.Status(http.StatusOK).JSON(JsonResponse{
//         Status:     true,
//         Message:    "Data retrieved successfully",
//         Data:       users,
//         StatusCode: http.StatusOK,
//     })
// }

func HandlerGetAllUser(c *fiber.Ctx) error {

    db, _ := initializers.ConnectToDB()

    query := `SELECT
     users.id,
     users.username,
     user_roles.name,
     users.api_key,
     users.client_id,
     countries.name,
     users.email,
     users.validation_token,
     users.mobile,
     users.referral_code,
     products.name,
     users.total_invitees,
     users.successful_referral,
     users.is_active
	 FROM users
	 LEFT JOIN user_roles ON user_roles.id = users.role_id
	 LEFT JOIN countries ON countries.code = users.country_code
	 LEFT JOIN products ON products.id = users.product_id
	 ORDER BY users.created_at DESC`

    row, err := db.Query(context.Background(), query)

    if err != nil {
		fmt.Println("---------error",err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "error in database in users",
        })
    }

    var users []Users
    defer row.Close()

    for row.Next() {
        var user Users
        if err = row.Scan(
            &user.ID,
            &user.Username,
            &user.RoleName,
            &user.ApiKey,
            &user.ClientID,
            &user.CountryName,
            &user.Email,
            &user.ValidationToken,
            &user.Mobile,
            &user.ReferralCode,
            &user.ProductName,
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

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message":users,
    })
}

// func HandlerGetAllUser(c *fiber.Ctx) error {
// 	db, _ := initializers.ConnectToDB()

// 	query := `SELECT
//      users.id,
//      users.username,
//      user_roles.name,
//      users.api_key,
//      users.client_id,
//      countries.name,
//      users.email,
//      users.validation_token,
//      users.mobile,
//      users.referral_code,
//      products.name,
//      users.total_invitees,
//      users.successful_referral,
//      users.is_active
// 	 FROM users
// 	 LEFT JOIN user_roles ON user_roles.id = users.role_id
// 	 LEFT JOIN countries ON countries.code = users.country_code
// 	 LEFT JOIN products ON products.id = users.product_id
// 	 ORDER BY users.created_at DESC`

// 	rows, err := db.Query(context.Background(), query)
// 	if err != nil {
// 		fmt.Println("Error executing query:", err)
// 		return c.Status(http.StatusBadRequest).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Error in database in users",
// 			StatusCode: http.StatusBadRequest,
// 		})
// 	}
// 	defer rows.Close()

// 	var users []Users
// 	for rows.Next() {
// 		var user Users
// 		if err = rows.Scan(
// 			&user.ID,
// 			&user.Username,
// 			&user.RoleName,
// 			&user.ApiKey,
// 			&user.ClientID,
// 			&user.CountryName,
// 			&user.Email,
// 			&user.ValidationToken,
// 			&user.Mobile,
// 			&user.ReferralCode,
// 			&user.ProductName,
// 			&user.TotalInvitees,
// 			&user.SuccessfulReferral,
// 			&user.IsActive,
// 		); err != nil {
// 			fmt.Println("Error scanning data:", err)
// 			return c.Status(http.StatusBadRequest).JSON(JsonResponse{
// 				Status:     false,
// 				Message:    "Error getting data from the database in users",
// 				StatusCode: http.StatusBadRequest,
// 			})
// 		}
// 		users = append(users, user)
// 	}

// 	if len(users) > 0 {
// 		return c.Status(http.StatusCreated).JSON(JsonResponse{
// 			Status:     true,
// 			Message:    "Data retrieved successfully",
// 			Data:       users,
// 			StatusCode: http.StatusCreated,
// 		})
// 	} else {
// 		return c.Status(http.StatusNotFound).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Data not found",
// 			Data:       nil,
// 			StatusCode: http.StatusNotFound,
// 		})
// 	}
// }

func HandlerGetOneUser(c *fiber.Ctx) error {
    db, err := initializers.ConnectToDB()
    if err != nil {
        return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
            "error": "unable to connect to the database",
        })
    }

    id := c.Params("id")

    query := `SELECT
    users.id,
     users.username,
     user_roles.name,
     users.api_key,
     users.client_id,
     countries.name,
     users.email,
     users.validation_token,
     users.mobile,
     users.referral_code,
     products.name,
     users.total_invitees,
     users.successful_referral,
     users.is_active
	FROM users
	LEFT JOIN user_roles ON user_roles.id = users.role_id
	LEFT JOIN countries ON countries.code = users.country_code
	LEFT JOIN products ON products.id = users.product_id
	where users.id = $1`

    row := db.QueryRow(context.Background(), query, id)

    var user Users
    err = row.Scan(
		&user.ID,
		&user.Username,
		&user.RoleName,
		&user.ApiKey,
		&user.ClientID,
		&user.CountryName,
		&user.Email,
		&user.ValidationToken,
		&user.Mobile,
		&user.ReferralCode,
		&user.ProductName,
		&user.TotalInvitees,
		&user.SuccessfulReferral,
		&user.IsActive,
    )

    if err != nil {
		fmt.Println("-----",err)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user": user,
    })
}

// func HandlerGetOneUser(c *fiber.Ctx) error {
//     db, err := initializers.ConnectToDB()
//     if err != nil {
//         return c.Status(fiber.StatusBadGateway).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Unable to connect to the database",
//             StatusCode: fiber.StatusBadGateway,
//         })
//     }

//     id := c.Params("id")

//     query := `SELECT
//     users.id,
//      users.username,
//      user_roles.name,
//      users.api_key,
//      users.client_id,
//      countries.name,
//      users.email,
//      users.validation_token,
//      users.mobile,
//      users.referral_code,
//      products.name,
//      users.total_invitees,
//      users.successful_referral,
//      users.is_active
// 	FROM users
// 	LEFT JOIN user_roles ON user_roles.id = users.role_id
// 	LEFT JOIN countries ON countries.code = users.country_code
// 	LEFT JOIN products ON products.id = users.product_id
// 	where users.id = $1`

//     row := db.QueryRow(context.Background(), query, id)

//     var user Users
//     err = row.Scan(
//         &user.ID,
//         &user.Username,
//         &user.RoleName,
//         &user.ApiKey,
//         &user.ClientID,
//         &user.CountryName,
//         &user.Email,
//         &user.ValidationToken,
//         &user.Mobile,
//         &user.ReferralCode,
//         &user.ProductName,
//         &user.TotalInvitees,
//         &user.SuccessfulReferral,
//         &user.IsActive,
//     )

//     if err != nil {
//         fmt.Println("Error:", err)
//         return c.Status(fiber.StatusNotFound).JSON(JsonResponse{
//             Status:     false,
//             Message:    "User not found",
//             StatusCode: fiber.StatusNotFound,
//         })
//     }

//     return c.Status(fiber.StatusOK).JSON(JsonResponse{
//         Status:     true,
//         Message:    "User retrieved successfully",
//         Data:       user,
//         StatusCode: fiber.StatusOK,
//     })
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

    query := `
        UPDATE users
        SET username = $1, role_id = $2, api_key = $3, client_id = $4, country_code = $5,
        email = $6, password = $7, validation_token = $8, mobile = $9, referral_code = $10,
        product_id = $11, total_invitees = $12, successful_referral = $13, is_active = $14
        WHERE id = $15
    `

    _, err = db.Exec(context.Background(), query,
        body.Username,
        body.RoleID,
        body.ApiKey,
        body.ClientID,
        body.CountryCode,
        body.Email,
        string(hash),
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

    c.Status(fiber.StatusAccepted).JSON(fiber.Map{
        "message": "User Updated Successfully",
    })

    return nil
}

// func HandleUpdateUser(c *fiber.Ctx) error {
//     db, _ := initializers.ConnectToDB()
//     id := c.Params("id")

//     var body models.User

//     if err := c.BodyParser(&body); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Failed to read body",
//             StatusCode: fiber.StatusBadRequest,
//         })
//     }

//     hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
//     if err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Unable to hash password",
//             StatusCode: fiber.StatusBadRequest,
//         })
//     }

//     query := `
//         UPDATE users
//         SET username = $1, role_id = $2, api_key = $3, client_id = $4, country_code = $5,
//         email = $6, password = $7, validation_token = $8, mobile = $9, referral_code = $10,
//         product_id = $11, total_invitees = $12, successful_referral = $13, is_active = $14
//         WHERE id = $15
//     `

//     _, err = db.Exec(context.Background(), query,
//         body.Username,
//         body.RoleID,
//         body.ApiKey,
//         body.ClientID,
//         body.CountryCode,
//         body.Email,
//         string(hash),
//         body.ValidationToken,
//         body.Mobile,
//         body.ReferralCode,
//         body.ProductID,
//         body.TotalInvitees,
//         body.SuccessfulReferral,
//         body.IsActive,
//         id,
//     )
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
//             Status:     false,
//             Message:    "Error updating user",
//             StatusCode: fiber.StatusInternalServerError,
//         })
//     }

//     return c.Status(fiber.StatusAccepted).JSON(JsonResponse{
//         Status:     true,
//         Message:    "User updated successfully",
//         StatusCode: fiber.StatusAccepted,
//     })
// }

func HandleDeleteUser(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	id := c.Params("id")

	query := `SELECT
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
	FROM users where id = $1`

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
		"message": "User Deleted Successfully",
	})

	return nil
}

// func HandleDeleteUser(c *fiber.Ctx) error {
// 	db, _ := initializers.ConnectToDB()

// 	id := c.Params("id")

// 	query := `DELETE FROM users WHERE id = $1 RETURNING
//     username,
//     role_id,
//     api_key,
//     client_id,
//     country_code,
//     email,
//     validation_token,
//     mobile,
//     referral_code,
//     product_id,
//     total_invitees,
//     successful_referral,
//     is_active`

// 	result, err := db.Exec(context.Background(), query, id)

// 	if err != nil {
// 		return c.Status(fiber.StatusBadGateway).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "Unable to execute the query",
// 			StatusCode: fiber.StatusBadGateway,
// 		})
// 	}

// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(JsonResponse{
// 			Status:     false,
// 			Message:    "User not found",
// 			StatusCode: fiber.StatusNotFound,
// 		})
// 	}

// 	return c.Status(fiber.StatusAccepted).JSON(JsonResponse{
// 		Status:     true,
// 		Message:    "User deleted successfully",
// 		StatusCode: fiber.StatusAccepted,
// 	})
// }


func HandleBulkDeleteUsers(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	var request struct {
		UserIDs []string `json:"user_ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(JsonResponse{
			Status:     false,
			Message:    "Failed to read request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	if len(request.UserIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(JsonResponse{
			Status:     false,
			Message:    "No user IDs provided for deletion",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := `DELETE FROM users WHERE id = ANY($1) RETURNING
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
    is_active`

	result, err := db.Exec(context.Background(), query, request.UserIDs)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(JsonResponse{
			Status:     false,
			Message:    "Unable to execute the query",
			StatusCode: fiber.StatusBadGateway,
		})
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(JsonResponse{
			Status:     false,
			Message:    "No users found for deletion",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(JsonResponse{
		Status:     true,
		Message:    "Users deleted successfully",
		StatusCode: fiber.StatusAccepted,
	})
}
