package auth

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	lock  sync.RWMutex
	jwks  *JWKS
	links *KeycloakLinks
}

func NewAuthentication(endpointsLink string) (*Authentication, error) {
	links, err := newKeycloakLinks(endpointsLink)
	if err != nil {
		return nil, err
	}

	jwks, err := newJWKS(links)
	if err != nil {
		return nil, err
	}

	return &Authentication{
		jwks:  jwks,
		links: links,
	}, nil
}

func (auth *Authentication) refreshJWKS() error {
	auth.lock.Lock()
	defer auth.lock.Unlock()

	newJWKS, err := newJWKS(auth.links)
	if err != nil {
		return err
	}

	auth.jwks = newJWKS
	return nil
}

func (auth *Authentication) validate(accessToken string) (*jwt.MapClaims, error) {
	auth.lock.RLock()
	defer auth.lock.RUnlock()

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		keyKid := token.Header["kid"].(string)
		pubKey, err := auth.jwks.getPublicKey(keyKid)
		if err != nil {
			return nil, err
		}

		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	switch v := token.Claims.(type) {
	case jwt.MapClaims:
		return &v, nil
	default:
		return nil, errors.New("claim type is not jwt.MapClaims")
	}
}
