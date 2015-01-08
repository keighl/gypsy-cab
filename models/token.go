package models

import (
  "time"
  "strings"
  "errors"
  r "github.com/dancannon/gorethink"
  u "gypsy/utils"
)

var tokenKey []byte = []byte("GLumCwoK89HhRykooYLZUBtk>=3m<:9q")

type Token struct {
  Errors []string `gorethink:"-" json:"errors,omitempty"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
  Id string `gorethink:"id,omitempty" json:"token"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at"`
  UpdatedAt time.Time `gorethink:"updated+at" json:"-"`
  UserId string `gorethink:"user_id" json:"-"`
  S3AccessId string `gorethink:"s3_access_id" json:"s3_access_id"`
  S3ClientSecretCrypted string `gorethink:"s3_client_secret" json:"-"`
  S3Bucket string `gorethink:"s3_bucket" json:"s3_bucket"`
  S3Region string `gorethink:"s3_region" json:"s3_region"`

}

type TokenAttrs struct {
  S3AccessId string `json:"s3_access_id" form:"s3_access_id"`
  S3ClientSecret string `json:"s3_client_secret" form:"s3_client_secret"`
  S3Bucket string `json:"s3_bucket" form:"s3_bucket"`
  S3Region string `json:"s3_region" form:"s3_region"`
}

//////////////////////////////
// TRANSACTIONS //////////////

func (x *Token) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  if (x.Id == "") {
    x.BeforeCreate()
    res, err := r.Table("tokens").Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    x.Id = res.GeneratedKeys[0]
  }

  x.BeforeUpdate()
  _, err := r.Table("tokens").Get(x.Id).Replace(x).RunWrite(DB)
  return err
}

func (x *Token) Delete() error {
  _, err := r.Table("tokens").Get(x.Id).Delete().RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *Token) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()
}

func (x *Token) BeforeUpdate() {
  x.UpdatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Token) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.Trimspace()
  x.ValidateS3Creds()
  return !x.HasErrors()
}

func (x *Token) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Token) Trimspace() {
  x.S3AccessId = strings.TrimSpace(x.S3AccessId)
  x.S3Bucket = strings.TrimSpace(x.S3Bucket)
  x.S3Region = strings.TrimSpace(x.S3Region)
}

func (x *Token) ValidateS3Creds() {

  if (x.S3AccessId == "") {
    x.Errors = append(x.Errors, "S3AccessId can't be blank")
    x.ErrorMap["S3AccessId"] = true
  }

  if (x.S3ClientSecretCrypted == "") {
    x.Errors = append(x.Errors, "S3ClientSecret can't be blank")
    x.ErrorMap["S3ClientSecret"] = true
  }

  if (x.S3Bucket == "") {
    x.Errors = append(x.Errors, "S3Bucket can't be blank")
    x.ErrorMap["S3Bucket"] = true
  }

  if (x.S3Region == "") {
    x.Errors = append(x.Errors, "S3Region can't be blank")
    x.ErrorMap["S3Region"] = true
  }
}

//////////////////////////////
// SECRET ////////////////////

func (x *Token) S3ClientSecret() string {
  if (x.S3ClientSecretCrypted == "") { return "" }
  return u.Decrypt(tokenKey, x.S3ClientSecretCrypted)
}

func (x *Token) SetSecret(secret string) {
  if (secret == "") { return }
  secret = strings.TrimSpace(secret)
  x.S3ClientSecretCrypted = u.Encrypt(tokenKey, secret)
}

//////////////////////////////
// OTHER /////////////////////

func (x *Token) UpdateFromAttrs(attrs TokenAttrs) {
  if (attrs.S3AccessId != "") { x.S3AccessId = attrs.S3AccessId }
  if (attrs.S3ClientSecret != "") { x.SetSecret(attrs.S3ClientSecret) }
  if (attrs.S3Bucket != "") { x.S3Bucket = attrs.S3Bucket }
  if (attrs.S3Region != "") { x.S3Region = attrs.S3Region }
}

func (x *TokenAttrs) Token() (*Token) {
  token := &Token{
    S3AccessId: x.S3AccessId,
    S3Bucket: x.S3Bucket,
    S3Region: x.S3Region,
  }
  token.SetSecret(x.S3ClientSecret)
  return token
}
