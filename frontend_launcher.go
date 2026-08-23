package main

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
  port := "3000"
  serverDir := "."

  if _, err := os.Stat("./index.html"); os.IsNotExist(err) {
    serveDir = "./frontend"
  }

  fmt.Printf("Frontend running on port 3000")

  fs := http.FileServer(http.Dir(serveDir))
  http.Handle("/", fs)

  err := http.ListenAndServe(":"+port, nil)
  if err != nil {
    fmt.Printf("Frontend server crashed: %v\n", err)
  }
}
