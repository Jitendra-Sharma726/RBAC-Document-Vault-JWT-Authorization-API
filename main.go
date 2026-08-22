package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"
  "time"

  "github.com/golang-jwt/jwt/v5"
)

//Secret Key used to sign and validate RBAC tokens.
//IMPORTANT: This is just a dummy key for educational purposes.
//Never hardcore real crytographic keys in your source code in production
var jwtKey = []byte("x8A9b2F4c1D7e3B6a0C5d8E2f9A1b4C3")

//User models our database record
type User struct {
  Username string
  Password string
  Role     string
}

//mockDB acts as our temporary in-memory database
var mockDB = map[string]User{
  "alice":   {Username: "alice", Password: "password123", Role: "admin"},
  "bob":     {Username: "bob", Password:   "password123", Role: "employee"},
  "charlie": {Username: "charlie", Password: "password123", Role: "guest"},
  }

type Credentials struct {
  Username string `json:"username"`
  Password string `json:"password"`

  //LoginHandler validates credentials and issues a role-embedded JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
     if r.Method != http.MethodPost {
       http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
       return
    }

  var creds Credentials
  err := json.NewDecoder(r.Body).Decode(&creds)
  if err != nil {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
  }

  user, exists := mockDB[creds.Username]
  if !exists || user.Password != creds.Password {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  // Construct JWT claims(username, role, exp)
  claims := jwt.MapClaims{
    "username":   user.Username,
    "role":       user.Role,
    "exp":        time.Now().Add(1 * time.Hour).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigninigMethodHS256, claims)

  //Sign the token using jwtKey
  tokenString, err  := token.SignedString(jwtKey)
  if err != nil {
     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
     return
  }

  //




















  













  
