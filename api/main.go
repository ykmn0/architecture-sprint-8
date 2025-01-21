package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	c0env "github.com/caarlos0/env"
)

const RetryCount = 10

type environments struct {
	KeycloakURL   string `env:"KEYCLOAK_URL"`
	KeycloakRealm string `env:"KEYCLOAK_REALM"`
}

var (
	keycloakURL = "http://localhost:8080/realms/reports-realm"
)

func init() {
	// Load environment variables
	env, err := getEnvironments()
	if err != nil {
		log.Fatalf("error loading environment variables: %v\n", err)
	}

	// Update URL and Realm values based on environment
	if env.KeycloakURL != "" {
		keycloakURL = env.KeycloakURL
	}
	if env.KeycloakRealm != "" {
		keycloakURL = fmt.Sprintf("%s/realms/%s", keycloakURL, env.KeycloakRealm)
	}
}

func main() {
	var publicKey *rsa.PublicKey

	// Attempt to get the public key with retries and increasing delay
	errRetry := RetryWithDelay(RetryCount, 500*time.Millisecond, func() (err error) {
		publicKey, err = FetchKeycloakPublicKey(keycloakURL)
		return err
	})

	if errRetry != nil {
		log.Fatalf("failed to get Keycloak public key after retries: %v\n", errRetry)
	}

	// Check if public key was received
	if publicKey == nil {
		log.Fatal("public key is nil")
	}

	// Handler for the /reports route
	http.HandleFunc("/reports", ReportsHandler(publicKey))

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Function to load environment variables
func getEnvironments() (*environments, error) {
	env := new(environments)
	err := c0env.Parse(env)
	return env, err
}

// Retry function with multiple attempts
func RetryWithDelay(attempts int, delay time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		log.Printf("attempt %d failed: %v; retrying in %v...", i+1, err, delay)
		time.Sleep(delay)
	}
	return fmt.Errorf("after %d attempts, last error: %w", attempts, fn())
}

// Example function to retrieve the Keycloak public key
func FetchKeycloakPublicKey(url string) (*rsa.PublicKey, error) {
	// Logic for retrieving the public key from Keycloak
	return nil, nil
}

// Example handler
func ReportsHandler(publicKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your request handling code here
		fmt.Fprintf(w, "You can use handler with the public key: %v", publicKey)
	}
}
