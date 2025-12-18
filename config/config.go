package config
import "os"

var FrontendURL string
var (
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	RefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

func LoadEnv() {
	FrontendURL = os.Getenv("FRONTEND_URL")
}