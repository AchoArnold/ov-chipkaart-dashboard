package jwt

import (
	"time"

	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	"github.com/palantir/stacktrace"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/cache"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

const (
	keyUserID = "user_id"
	keyExp    = "exp"
)

var (
	// ErrTokenBlacklisted is thrown when a jwt token is blacklisted
	ErrTokenBlacklisted = errors.New("token has been blacklisted")
)

var (
	// ErrCodeInvalidExpiryDate is the error code when claims[keyExp] is invalid
	ErrCodeInvalidExpiryDate = stacktrace.ErrorCode(1)
)

// Service is a new instance of the JWT service
type Service struct {
	secret      []byte
	cache       cache.Cache
	sessionDays int
}

// NewService creates a new instance of the JWT service
func NewService(secret string, cache cache.Cache, sessionDays int) Service {
	return Service{
		secret:      []byte(secret),
		cache:       cache,
		sessionDays: sessionDays,
	}
}

//GenerateTokenForUserID generates a jwt token and assign a email to it's claims and return it
func (service Service) GenerateTokenForUserID(UserID id.ID) (result string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims[keyExp] = time.Now().UTC().AddDate(0, 0, service.sessionDays).String()
	claims["nbf"] = time.Now().UTC().Unix()
	claims["iat"] = time.Now().UTC().Unix()
	claims[keyUserID] = UserID.String()

	result, err = token.SignedString(service.secret)
	if err != nil {
		return result, err
	}

	return result, nil
}

// IsValid checks if a token is valid
func (service Service) IsValid(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return false
	}

	_, err = service.cache.Get(tokenString)
	if err == nil {
		return false
	}

	return true
}

// InvalidateToken invalidates a jwt token
func (service Service) InvalidateToken(tokenString string) (err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil
	}

	expiryDate, ok := claims[keyExp].(string)
	if !ok {
		return stacktrace.NewErrorWithCode(ErrCodeInvalidExpiryDate, "cannot get expiry date from claim %v", claims[keyExp])
	}

	expiryTime, err := time.Parse(internalTime.StringFormat, expiryDate)
	if err != nil {
		return stacktrace.NewErrorWithCode(ErrCodeInvalidExpiryDate, "cannot parse claim expiry date %s with time format %s", expiryDate, internalTime.StringFormat)
	}

	return service.cache.Set(tokenString, "", expiryTime.Sub(time.Now().UTC()))
}

//GetUserIDFromToken parses a jwt token and returns the user ID
func (service Service) GetUserIDFromToken(tokenString string) (userID id.ID, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return userID, err
	}

	_, err = service.cache.Get(tokenString)
	if err == nil {
		return userID, ErrTokenBlacklisted
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, err
	}

	return id.FromString(claims[keyUserID].(string))
}
