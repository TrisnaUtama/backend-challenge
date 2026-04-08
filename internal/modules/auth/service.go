package auth

import (
	"crypto/rand"
	"encoding/hex"
)

type Service interface {
	GenerateToken() (TokenResponse, error)
	ValidateToken(t string) bool
}

type service struct {
	token string
}

func NewService() Service {
	return &service{}
}

func (s *service) GenerateToken() (TokenResponse, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return TokenResponse{}, err
	}

	token := hex.EncodeToString(b)
	s.token = token

	return TokenResponse{
		Token: token,
	}, nil
}

func (s *service) ValidateToken(t string) bool {
	return t == s.token
}
