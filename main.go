package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// User Dto
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

var jwtSecret = []byte("thepolyglotdeveloper")

func main() {
	r := mux.NewRouter()
	r.Handle("/authenticate", AuthenticationHandler).Methods("POST")
	r.Handle("/products", ValidateMiddleware(ProductsHandler)).Methods("GET")

	_ = http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}

// AuthenticationHandler is the handler to authenticate user and create new token
var AuthenticationHandler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)

	fmt.Println(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
})

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")

		if authorizationHeader == "" {
			_ = json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		} else {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return jwtSecret, nil
				})
				if error != nil {
					_ = json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					_ = json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			} else {
				_ = json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
			}
		}
	})
}

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
	var user User
	_ = mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	_ = json.NewEncoder(w).Encode(user)
})

// NotImplemented handler
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Not Implemented"))
})
