package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
)

// Validate is a struct that holds the necessary information to validate a jwt
type Validate struct {
	logger          *slog.Logger
	clientId        string
	authorizedUsers map[string]bool
	jwksUrl         string
	jwksSet         jwk.Set
}

// NewValidate creates a new Validate struct. This is meant to be used as a middleware to validate jwt tokens.
func NewValidate(logger *slog.Logger, clientId string, authorizedUsers map[string]bool, jwksUrl string) *Validate {
	set, err := getJwksSet(jwksUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get jwks set, error: %v", err))
	}

	return &Validate{
		logger:          logger,
		clientId:        clientId,
		authorizedUsers: authorizedUsers,
		jwksUrl:         jwksUrl,
		jwksSet:         set,
	}
}

func (v *Validate) IsValidJwt(r *http.Request) (bool, error) {
	authToken := r.Header.Get("Authorization")
	tokenString := strings.Replace(authToken, "Bearer ", "", 1)
	if tokenString == "" {
		v.logger.Warn("auth token not present")
		return false, fmt.Errorf("auth token not present")
	}

	token, err := v.parseAndValidateToken(tokenString)
	if err != nil {
		v.logger.Warn(fmt.Sprintf("auth failed to parse and validate token, error: %v", err.Error()))
		return false, err
	}

	if !token.Valid {
		v.logger.Warn("auth token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	if v.isValidEmail(claims) && v.isValidClient(claims) {
		return true, nil
	}

	return false, fmt.Errorf("auth invalid token claims")
}

func (v *Validate) parseAndValidateToken(tokenString string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error
	token, err = jwt.Parse(tokenString, v.getKey, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		v.logger.Warn(fmt.Sprintf("auth failed to parse and validate token in first attempt, error: %v", err.Error()))
		v.logger.Warn("auth retrying to validate token after updating jwks cache")
		set, err := getJwksSet(v.jwksUrl)
		if err != nil {
			return nil, err
		}
		v.jwksSet = set
		v.logger.Info("auth successfully updated jwks cache")
		token, err = jwt.Parse(tokenString, v.getKey, jwt.WithValidMethods([]string{"RS256"}))
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func (v *Validate) getKey(token *jwt.Token) (interface{}, error) {
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		v.logger.Warn("auth expecting JWT header to have string kid")
		return nil, fmt.Errorf("auth expecting JWT header to have string kid")
	}

	key, err := v.lookupKeyId(keyID)
	if err != nil {
		return nil, err
	}

	var pubKey rsa.PublicKey
	err = key.Raw(&pubKey)
	if err != nil {
		v.logger.Error("auth failed to create rsa.PublicKey key from jwks set")
		return nil, err
	}
	return &pubKey, nil
}

func (v *Validate) lookupKeyId(keyId string) (jwk.Key, error) {
	key, ok := v.jwksSet.LookupKeyID(keyId)
	if !ok { // key not found in jwksSet; refetch from jwks end-point
		v.logger.Info("auth key not found in jwksSet, updating jwks cache")
		set, err := getJwksSet(v.jwksUrl)
		if err != nil {
			v.logger.Error(fmt.Sprintf("auth failed to fetch from jwks endpoint, error: %v", err.Error()))
			return nil, err
		}
		v.jwksSet = set
		v.logger.Info("auth successfully updated jwks cache")

		// now again check for key in key set
		jkey, ok := v.jwksSet.LookupKeyID(keyId)
		if !ok {
			v.logger.Warn("auth kid not found in jwks set")
			return nil, err
		}
		return jkey, nil
	}

	return key, nil
}

func (v *Validate) isValidEmail(claims jwt.MapClaims) bool {
	var email string
	if userEmail, ok := claims["email"]; ok {
		email = strings.ToLower(userEmail.(string))
	} else {
		v.logger.Warn("auth token validation failure, email claim not found in jwt token")
		return false
	}

	if _, ok := v.authorizedUsers[email]; ok {
		return true
	} else {
		v.logger.Warn("auth token validation failure, unauthorized email id")
		return false
	}
}

func (v *Validate) isValidClient(claims jwt.MapClaims) bool {
	incomingClientid, ok := claims["client_id"]
	if !ok {
		v.logger.Warn("auth token validation failure, client_id claim not found in jwt token")
		return false
	}

	if v.clientId == incomingClientid {
		return true
	} else {
		v.logger.Warn("auth token validation failure, invalid sso client id")
		return false
	}
}

func getJwksSet(jwksUrl string) (jwk.Set, error) {
	// TODO: add retry logic
	set, err := jwk.Fetch(context.Background(), jwksUrl)
	if err != nil {
		return nil, err
	}

	return set, nil
}
