package main

import (
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const ProtheticUser = "prothetic_user"

var (
	headers = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Token",
	}
)

// Handler returns the handler function for the /reports route.
func Handler(publicKey *rsa.PublicKey) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Enable CORS headers for the response
		enableCors(writer)

		// If the method is OPTIONS, respond with a status OK
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}

		// Extract token from the request
		reqToken := GetToken(request)

		// Validate the extracted token
		jwtToken, errValidateToken := ValidateToken(reqToken, publicKey)
		if errValidateToken != nil {
			log.Printf("error: %v\n", errValidateToken)
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write([]byte(errValidateToken.Error()))
			return
		}

		// Extract roles from the JWT token
		roles := RolesClaims{
			Claims: jwtToken.Claims.(jwt.MapClaims),
		}.ToMap().
			ToSlice().
			ToRoles()

		// Check if the user has the "prothetic_user" role
		if _, ok := roles[ProtheticUser]; !ok {
			log.Printf("error: Forbidden\n")
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		// If everything is valid, return OK status
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("Ok"))
	}
}

// enableCors sets the necessary CORS headers for the response
func enableCors(w http.ResponseWriter) {
	for header, value := range headers {
		w.Header().Set(header, value)
	}
}
