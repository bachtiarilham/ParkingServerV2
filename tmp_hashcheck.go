package main
import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)
func main(){
  err:=bcrypt.CompareHashAndPassword([]byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"), []byte("password"))
  if err != nil { panic(err) }
  fmt.Print("ok")
}
