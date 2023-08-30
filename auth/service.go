package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(UserID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

// SECRET KEY
var SECRET_KEY = []byte("POLYMATIC_s3cr3T_K3Y")

func (s *jwtService) GenerateToken(UserID int) (string, error) {

	// Payload: ambil userID untuk di include-kan pada generate Token (Bisa ambil ID, Nama and etc..)
	payload := jwt.MapClaims{}
	payload["user_id"] = UserID

	// Set Header JWT: Menentukan jenis algoritma & Token Type
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Token yang sudah dibuat wajib harus di tanda-tangani (Signature)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
