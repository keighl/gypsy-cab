package models

import (
  "testing"
  u "gypsy/utils"
)

func NewTestToken() *Token {
  return &Token{
    S3AccessId: "XXX",
    S3ClientSecretCrypted: "XXX",
    S3Bucket: "XXX",
    S3Region: "XXX",
  }
}

//////////////////////////////
// VALIDATIONS ///////////////

func Test_Token_RequiresS3Creds(t *testing.T) {

  x := NewTestToken()

  x.S3AccessId = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["S3AccessId"], true)

  x.S3ClientSecretCrypted = ""
  err = Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["S3ClientSecret"], true)

  x.S3Bucket = ""
  err = Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["S3Bucket"], true)

  x.S3Region = ""
  err = Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["S3Region"], true)

  x = NewTestToken()
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["S3AccessId"], false)
  expect(t, x.ErrorMap["S3ClientSecret"], false)
  expect(t, x.ErrorMap["S3Bucket"], false)
  expect(t, x.ErrorMap["S3Region"], false)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Token_UpdateAttrs(t *testing.T) {
  var tokenKey []byte = []byte(Config.TokenEncryptionKey)

  obj := &Token{
    S3AccessId: "cheese",
    S3ClientSecretCrypted: "cheese",
    S3Bucket: "cheese",
    S3Region: "cheese",
  }
  attrs := TokenAttrs{
    S3AccessId: "xcheese",
    S3ClientSecret: "xcheese",
    S3Bucket: "xcheese",
    S3Region: "xcheese",
  }
  obj.UpdateFromAttrs(attrs)
  targetByHand := &Token{
    S3ClientSecretCrypted: u.Encrypt(tokenKey, attrs.S3ClientSecret),
    S3AccessId: attrs.S3AccessId,
    S3Bucket: attrs.S3Bucket,
    S3Region: attrs.S3Region,
  }

  expect(t, u.Decrypt(tokenKey, targetByHand.S3ClientSecretCrypted), attrs.S3ClientSecret)
  expect(t, obj.S3AccessId, targetByHand.S3AccessId)
  expect(t, obj.S3Bucket, targetByHand.S3Bucket)
  expect(t, obj.S3Region, targetByHand.S3Region)
}

func Test_TokenAttrs_Token(t *testing.T) {
  var tokenKey []byte = []byte(Config.TokenEncryptionKey)

  obj := &TokenAttrs{
    S3AccessId: "cheese",
    S3ClientSecret: "cheese",
    S3Bucket: "cheese",
    S3Region: "cheese",
  }
  targetByMethod := obj.Token()
  targetByHand := &Token{
    S3ClientSecretCrypted: u.Encrypt(tokenKey, obj.S3ClientSecret),
    S3AccessId: obj.S3AccessId,
    S3Bucket: obj.S3Bucket,
    S3Region: obj.S3Region,
  }

  expect(t, u.Decrypt(tokenKey, targetByHand.S3ClientSecretCrypted), targetByMethod.S3ClientSecret())
  expect(t, targetByMethod.S3AccessId, targetByHand.S3AccessId)
  expect(t, targetByMethod.S3Bucket, targetByHand.S3Bucket)
  expect(t, targetByMethod.S3Region, targetByHand.S3Region)
}



