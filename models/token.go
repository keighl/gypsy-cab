package models

import (
  "strings"
  u "gypsy/utils"
)

var tokenKey []byte = []byte("GLumCwoK89HhRykooYLZUBtk>=3m<:9q")

type Token struct {
  Record
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

func (x *Token) Table() string {
  return "tokens"
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Token) Validate() {
  x.Record.Validate()
  x.Trimspace()
  x.ValidateS3Creds()
}

func (x *Token) Trimspace() {
  x.S3AccessId = strings.TrimSpace(x.S3AccessId)
  x.S3Bucket = strings.TrimSpace(x.S3Bucket)
  x.S3Region = strings.TrimSpace(x.S3Region)
}

func (x *Token) ValidateS3Creds() {

  if (x.S3AccessId == "") {
    x.ErrorOn("S3AccessId", "S3AccessId can't be blank")
  }

  if (x.S3ClientSecretCrypted == "") {
    x.ErrorOn("S3ClientSecret", "S3ClientSecret can't be blank")
  }

  if (x.S3Bucket == "") {
    x.ErrorOn("S3Bucket", "S3Bucket can't be blank")
  }

  if (x.S3Region == "") {
    x.ErrorOn("S3Region", "S3Region can't be blank")
  }
}

//////////////////////////////
// SECRET ////////////////////

func (x *Token) S3ClientSecret() string {
  return u.Decrypt(tokenKey, x.S3ClientSecretCrypted)
}

func (x *Token) SetSecret(secret string) {
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
