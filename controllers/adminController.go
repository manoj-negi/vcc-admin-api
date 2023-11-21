package controllers

import (
	"context"
	"errors"
	"time"
	"github.com/keshav/fiber/auth"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/models"
	
)

type AdminLogin struct{
	Email   string   `json:"email"`
	Password    string  `json:"password"`
}

func Login(username, password string) (*models.User, string, time.Time, error) {
    
    db, err := initializers.ConnectToDB()
    if err != nil {
        return nil, "", time.Time{}, err
    }

    var user models.User
    err = db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1", username).Scan(
        &user.ID,
        &user.Username,
        &user.RoleID,
        &user.ApiKey,
        &user.ClientID,
        &user.CountryCode,
        &user.Email,
        &user.Password,
        &user.ValidationToken,
        &user.Mobile,
        &user.ReferralCode,
        &user.ProductID,
        &user.TotalInvitees,
        &user.SuccessfulReferral,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, "", time.Time{}, errors.New("Invalid email or password")
        }
        return nil, "", time.Time{}, err
    }
	

    hashedPassword := auth.Sha256Hash(password)

    if user.Password != hashedPassword {

        return nil, "", time.Time{}, errors.New("Invalid email or password")
    }

    accessToken, expirationTime, err := auth.GenerateJWT(&user)
    if err != nil {
        return nil, "", time.Time{}, err
    }

    return &user, accessToken, expirationTime, nil
}

func HandlerAdminLogin(c *fiber.Ctx) error {
    var login AdminLogin
    if err := c.BodyParser(&login); err != nil {
        response := JsonResponse{
            Status:     false,
            Message:    "Invalid request body",
            StatusCode: fiber.StatusBadRequest,
        }
        return c.Status(fiber.StatusBadRequest).JSON(response)
    }

    user, accessToken, expirationTime, err := Login(login.Email, login.Password)
    if err != nil {
        response := JsonResponse{
            Status:     false,
            Message:    err.Error(),
            StatusCode: fiber.StatusUnauthorized,
        }
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }

    apiResponse := map[string]interface{}{
        "token":      accessToken,
        "userEmail":  user.Email,
        "expireDate": expirationTime.Format("2006-01-02"),
    }

    response := JsonResponse{
        Status:     true,
        Message:    "Authenticated successfully",
        Data:       apiResponse,
        StatusCode: fiber.StatusOK,
    }

    return c.Status(fiber.StatusOK).JSON(response)
}
























