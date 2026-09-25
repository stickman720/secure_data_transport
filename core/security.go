package core



import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"fmt"
	"errors"
)




type Security struct {
	JWTSecret  string
	HashSecret string	
}



type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"username"`
	Role string `json:"role"`
}




func (s *Security) CreateJWT(username string, role string, ttl time.Duration) (string, error) {
	now := time.Now()

	c := jwt.MapClaims{
		"username": username,
		"role":     role,
		"iat":      now.Unix(),
		"exp":      now.Add(ttl).Unix(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	toks, err := tok.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}

	return toks, nil
}

func (s *Security)ParseToken( tokenStr string) (*Claims, error) {
	secret := []byte(s.JWTSecret)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("myapp"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
