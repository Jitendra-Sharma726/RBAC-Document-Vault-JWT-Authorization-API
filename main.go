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

  //Return the generated token as JSON
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{
       "token": tokenString,
    })
}

// rbacMiddleware intercepts requests, validates JWTs, and enforces role boundaries
func rbacMiddleware(allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

      //Extract the token from the Authorization header
      authHeader := r.Header.Get("Authorization")
      if authHeader == "" {
        http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
        return
      }
      tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

      //Parse and validate the JWT signature
      token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        _, ok := token.Method(*jwt.SigningMethodHMAC);

        if !ok {
             return nil, fmt.Errorf("unexpected signing method")
        }
        return jwtKey, nil
        })

       if err != nil || !token.Valid {
          http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
         return
        }

      //Extract the user's role from the token claims
      claims, ok := token.Claims.(jwt.MapClaims)
      if !ok {
          http.Error(w, "Invalid token claims", http.StatusUnauthorized)
          return
      }

      userRole, ok := claims["role"].(string)
      if !ok {
           http.Error(w, "Role not found in token", http.StatusUnauthorized)
           return
      }

      //Deny access if the user's role is not in allowedRoles
      hasPermission := false
      for _, role := range allowedRoles {
          if userRole == role {
            hasPermission = true
            break
            }
        }

    if !hasPermission {
         http.Error(w, "Forbidden: You do not have the required role to access this document.", http.StatusForbidden)
         return
    }

    // Grant access to the requested document
    next(w, r)
  }
}

  func main() {
    fmt.Println("=== RBAC Document Vault ===")
    fmt.Println("Listening on port 8080...")

    mux := http.NewServeMux()

    //Public login endpoint
    mux.HandleFunc("/login", LoginHandler)

    //Wire up the router with RBAC middleware
    mux.HandleFunc("/docs/public", rbacMiddleware([]string{"admin", "employee", "guest"}, PublicDocHandler))
    











        

      



















  













  
