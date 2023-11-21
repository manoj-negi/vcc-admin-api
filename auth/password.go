package auth
import(
	"fmt"
	"time"
	"crypto/sha256"
    "encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"github.com/keshav/fiber/models"
)

var (
	accessTokenSecret = []byte("KXsMPri4PLlFlGcqU0f4P9y2s0aIOos9")
)

func GenerateJWT(user *models.User) (string, time.Time, error) {
	expirationTime := time.Now().AddDate(0, 2, 0)

	claims := &models.Claims{
		UserID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(accessTokenSecret)
	if err != nil {
		fmt.Println(err)
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func Sha256Hash(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

