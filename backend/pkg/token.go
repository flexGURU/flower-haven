package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID      uuid.UUID `json:"id"`
	UserID  uint32    `json:"user_id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	IsAdmin bool      `json:"is_admin"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	secretKey   string
	tokenIssuer string
}

func NewJWTMaker(secretKey, tokenIssuer string) JWTMaker {
	return JWTMaker{secretKey: secretKey, tokenIssuer: tokenIssuer}
}

func (maker *JWTMaker) CreateToken(userID uint32, name string, email string, duration time.Duration, isAdmin bool) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to create uuid: %v", err)
	}

	claims := Payload{
		id,
		userID,
		name,
		email,
		isAdmin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    maker.tokenIssuer,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to create token: %v", err)
	}

	return token, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, Errorf(INTERNAL_ERROR, "unexpected signing method")
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token: %v", err)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token is invalid")
	}

	if payload.RegisteredClaims.Issuer != maker.tokenIssuer {
		return nil, Errorf(INTERNAL_ERROR, "invalid issuer")
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, Errorf(INTERNAL_ERROR, "token is expired")
	}

	return payload, nil
}
