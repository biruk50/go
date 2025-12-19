package Infrastructure

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Generate(userID, username, role string, ttl time.Duration) (string, error)
	Validate(tokenStr string) (*jwt.RegisteredClaims, map[string]interface{}, error)
}

type jwtService struct {
	secret string
}

func GetJWTSecret() string {
	return GetEnv("JWT_SECRET", "")
}

func NewJWTService() JWTService {
	secret := GetJWTSecret()
	return &jwtService{secret: secret}
}

func (j *jwtService) Generate(userID, username, role string, ttl time.Duration) (string, error) {
	if j.secret == "" {
		j.secret = os.Getenv("JWT_SECRET")
	}
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(ttl).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) Validate(tokenStr string) (*jwt.RegisteredClaims, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}
	claimsMap := map[string]interface{}{}
	if mc, ok := token.Claims.(jwt.MapClaims); ok {
		for k, v := range mc {
			if k == "exp" || k == "iat" || k == "sub" {
			}
			claimsMap[k] = v
		}
	}
	// note: we return nil RegisteredClaims because parsing into RegisteredClaims would require re-parsing
	return nil, claimsMap, nil
}
