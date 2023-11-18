package models

import(
	"time"
)

// type User struct {
// 	Id             int    `json:"id"`
// 	First_name     string `json:"first_name"`
// 	Last_name      string `json:"last_name"`
// 	Email          string `json:"email"`
// 	Password       string `json:"password"`
// 	Age            int    `json:"age"`
// 	Phone_no       int    `json:"phone_no"`
// 	Secret_code    string `json:"secret_code"`
// 	Role_id        int    `json:"role_id"`
//}
type User struct {
	ID                 int      `json:"id"`
	Username           string   `json:"username"`
	RoleID             int      `json:"role_id"`
	ApiKey             string   `json:"api_key"`
	ClientID           string   `json:"client_id"`
	CountryCode        int      `json:"country_code"`
	Email              string   `json:"email"`
	Password           string   `json:"password"`
	ValidationToken    string   `json:"validation_token"`
	Mobile             string   `json:"mobile"`
	ReferralCode       string   `json:"referral_code"`
	ProductID          int      `json:"product_id"`
	TotalInvitees      int      `json:"total_invitees"`
	SuccessfulReferral int      `json:"successful_referral"`
	IsActive           int      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`  
	UpdatedAt          time.Time `json:"updated_at"`
}

type JsonResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

// type Role struct{
// 	Id int `json:"id"`
// 	Role string `json:"role_name"`
// }
// type Permissions struct{
// 	Id int `json:"id"`
// 	Permission string `json:"permission_name"`
// }
// type Role_Permission struct {
// 	Id int `json:"id"`
// 	Role_id int `json:"role_id"`
// 	Permission_id int `json:"permission_id"`
// }
// type Images struct{
// 	Image_url string `gorm:"image_url"`
// }