package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const PublicKey = "public_key"

var ErrNoPublicKey = errors.New("no public_key found")

// GetKeycloakPublicKey retrieves the public RSA key from Keycloak using its URL.
// Keycloak uses RS256 by default for signing JWT tokens.
func GetKeycloakPublicKey(keycloakURL string) (*rsa.PublicKey, error) {
	client := http.Client{}
	resp, errGet := client.Get(keycloakURL)
	if errGet != nil {
		return nil, errGet
	}
	bytes, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	// Unmarshal the response into a map
	issuer := make(map[string]interface{})
	errUnmarshal := json.Unmarshal(bytes, &issuer)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	// Retrieve the public key from the response map
	var base64EncodedPublicKey string
	if b, ok0 := issuer[PublicKey]; !ok0 {
		return nil, ErrNoPublicKey
	} else {
		base64EncodedPublicKey = b.(string)
	}

	// Parse and return the RSA public key
	return parseKeycloakRSAPublicKey(base64EncodedPublicKey)
}

// parseKeycloakRSAPublicKey decodes the base64-encoded public key and parses it into an *rsa.PublicKey.
func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}

	// Parse the public key using the X.509 standard
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}

	// Ensure the parsed key is an RSA public key
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
