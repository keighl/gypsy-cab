package utils

import (
  "os"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/cors"
  _ "github.com/go-sql-driver/mysql"
  r "github.com/dancannon/gorethink"
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/base64"
  "fmt"
  "io"
)

type Configuration struct {
  AppName string
  BaseURL string
  RethinkHost string
  RethinkDatabase string
  ServerLoggingEnabled bool
  MandrillAPIKey string
}

func ConfigForFile(confFile string) *Configuration {
  file, err := os.Open(confFile)
  if (err != nil) { panic(err) }
  decoder := json.NewDecoder(file)
  c := &Configuration{}
  err = decoder.Decode(c)
  if (err != nil) { panic(err) }
  file.Close()
  return c
}

func RethinkSession(conf *Configuration) *r.Session {
  session, _ := r.Connect(r.ConnectOpts{
    Address:  conf.RethinkHost,
    Database: conf.RethinkDatabase,
  })
  return session
}

func MartiniServer(logginEnabled bool) (*martini.ClassicMartini) {
  router := martini.NewRouter()
  server := martini.New()
  if (logginEnabled) { server.Use(martini.Logger()) }
  server.Use(martini.Recovery())
  server.MapTo(router, (*martini.Routes)(nil))
  server.Action(router.Handle)
  s := &martini.ClassicMartini{server, router}
  s.Use(render.Renderer())
  s.Use(cors.Allow(&cors.Options{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
    AllowHeaders: []string{"*", "x-requested-with", "Content-Type", "If-Modified-Since", "If-None-Match", "X-API-TOKEN"},
    ExposeHeaders: []string{"Content-Length"},
  }))
  return s
}

// https://gist.github.com/manishtpatel/8222606
// encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
  // key := []byte(keyText)
  plaintext := []byte(text)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  ciphertext := make([]byte, aes.BlockSize+len(plaintext))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  // convert to base64
  return base64.URLEncoding.EncodeToString(ciphertext)
}

// https://gist.github.com/manishtpatel/8222606
// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
  ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  if len(ciphertext) < aes.BlockSize {
    panic("ciphertext too short")
  }
  iv := ciphertext[:aes.BlockSize]
  ciphertext = ciphertext[aes.BlockSize:]

  stream := cipher.NewCFBDecrypter(block, iv)

  // XORKeyStream can work in-place if the two arguments are the same.
  stream.XORKeyStream(ciphertext, ciphertext)

  return fmt.Sprintf("%s", ciphertext)
}