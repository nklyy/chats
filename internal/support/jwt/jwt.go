package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	Uid  string `json:"uid"`
	jwt.StandardClaims
}

type CachedTokens struct {
	AccessUid  string `json:"access"`
	RefreshUid string `json:"refresh"`
}

//go:generate mockgen -source=jwt.go -destination=mocks/jwt_mock.go
type Service interface {
	CreateTokens(ctx context.Context, id, role string) (*string, *string, error)
	ParseToken(token string, isAccess bool) (*Payload, error)
	VerifyToken(ctx context.Context, payload *Payload, isAccess bool) error
	DeleteTokens(ctx context.Context, payload *Payload) error
	ExtendExpire(ctx context.Context, payload *Payload) error
}

type service struct {
	secretKeyAccess  string
	expiryAccess     int
	secretKeyRefresh string
	expiryRefresh    int
	autoLogout       int
	redisClient      *redis.Client
}

func NewJwtService(secretKeyAccess string,
	expiryAccess *int,
	secretKeyRefresh string,
	expiryRefresh *int,
	autoLogout *int,
	redisClient *redis.Client) (Service, error) {
	if secretKeyAccess == "" {
		return nil, errors.New("invalid jwt access secret key")
	}
	if expiryAccess == nil {
		return nil, errors.New("invalid jwt expiry access")
	}
	if secretKeyRefresh == "" {
		return nil, errors.New("invalid jwt refresh secret key")
	}
	if expiryRefresh == nil {
		return nil, errors.New("invalid jwt expiry refresh")
	}
	if autoLogout == nil {
		return nil, errors.New("invalid jwt auto logout")
	}
	if redisClient == nil {
		return nil, errors.New("invalid redis client")
	}
	return &service{
		secretKeyAccess:  secretKeyAccess,
		expiryAccess:     *expiryAccess,
		secretKeyRefresh: secretKeyRefresh,
		expiryRefresh:    *expiryRefresh,
		autoLogout:       *autoLogout,
		redisClient:      redisClient}, nil
}

func (s *service) CreateTokens(ctx context.Context, id, role string) (*string, *string, error) {
	// sign access token
	accessUid := uuid.New().String()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Payload{
		Id:   id,
		Role: role,
		Uid:  accessUid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.expiryAccess)).Unix(),
		},
	})

	signedAccessToken, err := accessToken.SignedString([]byte(s.secretKeyAccess))
	if err != nil {
		return nil, nil, ErrToken
	}

	// sign refresh token
	refreshUid := uuid.New().String()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Payload{
		Id:   id,
		Role: role,
		Uid:  refreshUid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.expiryRefresh)).Unix(),
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(s.secretKeyRefresh))
	if err != nil {
		return nil, nil, ErrToken
	}

	cacheJson, err := json.Marshal(CachedTokens{
		AccessUid:  accessUid,
		RefreshUid: refreshUid,
	})
	if err != nil {
		return nil, nil, ErrFailedCreateCache
	}

	err = s.redisClient.Set(ctx, fmt.Sprintf("token-%v", id), string(cacheJson), time.Minute*time.Duration(s.autoLogout)).Err()
	if err != nil {
		return nil, nil, ErrFailedCreateTokens
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func (s *service) ParseToken(token string, isAccess bool) (*Payload, error) {
	var secret string
	switch isAccess {
	case true:
		secret = s.secretKeyAccess
	case false:
		secret = s.secretKeyRefresh
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrToken
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return nil, ErrTokenInvalidOrExpire
		}
		return nil, ErrToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrToken
	}

	return payload, nil
}

func (s *service) VerifyToken(ctx context.Context, payload *Payload, isAccess bool) error {
	cacheJSON, err := s.redisClient.Get(ctx, fmt.Sprintf("token-%v", payload.Id)).Result()
	if err != nil {
		return ErrToken
	}

	cachedTokens := new(CachedTokens)
	err = json.Unmarshal([]byte(cacheJSON), cachedTokens)
	if err != nil {
		return err
	}

	var tokenUid string
	switch isAccess {
	case true:
		tokenUid = cachedTokens.AccessUid
	case false:
		tokenUid = cachedTokens.RefreshUid
	}

	if err != nil || tokenUid != payload.Uid {
		return ErrNotFound
	}

	return nil
}

func (s *service) DeleteTokens(ctx context.Context, payload *Payload) error {
	err := s.redisClient.Del(ctx, fmt.Sprintf("token-%v", payload.Id)).Err()
	if err != nil {
		return ErrFailedDeleteToken
	}

	return nil
}

func (s *service) ExtendExpire(ctx context.Context, payload *Payload) error {
	err := s.redisClient.Expire(ctx, fmt.Sprintf("token-%v", payload.Id), time.Minute*time.Duration(s.autoLogout)).Err()
	if err != nil {
		return ErrFailedExtendToken
	}

	return nil
}
