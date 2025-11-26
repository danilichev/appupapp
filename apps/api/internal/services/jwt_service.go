package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"apps/api/internal/api"
	"apps/api/internal/config"
)

type JwtClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type JWTService struct {
	refreshExpiration time.Duration
	refreshKey        string
	secretExpiration  time.Duration
	secretKey         string
}

func NewJWTService(
	config *config.JwtConfig,
) *JWTService {
	return &JWTService{
		secretKey: config.SecretKey,
		secretExpiration: time.Duration(
			config.SecretExpirationMinutes,
		) * time.Minute,
		refreshKey: config.RefreshKey,
		refreshExpiration: time.Duration(
			config.RefreshExpirationMinutes,
		) * time.Minute,
	}
}

func (s *JWTService) GenerateAuthToken(
	userId string,
) (*api.AuthToken, error) {
	accessToken, err := GenerateJwtToken(
		NewJwtClaims(userId, s.secretExpiration),
		s.secretKey,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateJwtToken(
		NewJwtClaims(userId, s.refreshExpiration),
		s.refreshKey,
	)
	if err != nil {
		return nil, err
	}

	return &api.AuthToken{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}

func (s *JWTService) ParseAccessToken(
	tokenString string) (*JwtClaims, error) {
	return ParseJwtToken(
		tokenString,
		s.secretKey,
	)
}

func (s *JWTService) ParseRefreshToken(
	tokenString string) (*JwtClaims, error) {
	return ParseJwtToken(
		tokenString,
		s.refreshKey,
	)
}

func NewJwtClaims(userId string, duration time.Duration) JwtClaims {
	claims := JwtClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return claims
}

func GenerateJwtToken(claims JwtClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJwtToken(tokenString string, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
