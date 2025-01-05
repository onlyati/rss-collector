package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

type KeycloakLinks struct {
	Issuer        string `json:"issuer"`
	TokenEndpoint string `json:"token_endpoint"`
	JWKSUri       string `json:"jwks_uri"`
}

func newKeycloakLinks(configLink string) (*KeycloakLinks, error) {
	resp, err := http.Get(configLink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var links KeycloakLinks
	if err := json.Unmarshal(body, &links); err != nil {
		return nil, err
	}

	return &links, nil
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func newJWKS(links *KeycloakLinks) (*JWKS, error) {
	resp, err := http.Get(links.JWKSUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func (jwks *JWKS) getPublicKey(kid string) (*rsa.PublicKey, error) {
	for _, key := range jwks.Keys {
		if kid == key.Kid {
			decodeN := decodeBase64URL(key.N)
			decodeE := decodeBase64URL(key.E)

			if decodeN == nil || decodeE == nil {
				return nil, errors.New("failed to decode key parameters")
			}

			pubKey := &rsa.PublicKey{
				N: decodeN,
				E: int(decodeE.Uint64()),
			}
			return pubKey, nil
		}
	}

	return nil, errors.New("public key not found for the given kid")
}

func decodeBase64URL(s string) *big.Int {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	return new(big.Int).SetBytes(bytes)
}
