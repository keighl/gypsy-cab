package models

import (
  "testing"
  "reflect"
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
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["S3AccessId"], true)

  x.S3ClientSecretCrypted = ""
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["S3ClientSecret"], true)

  x.S3Bucket = ""
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["S3Bucket"], true)

  x.S3Region = ""
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["S3Region"], true)

  x = NewTestToken()
  x.Save()
  expect(t, x.Validate(), true)

}

//////////////////////////////
// TRANSACTIONS //////////////

func Test_Token_Create_Success(t *testing.T) {

  x := NewTestToken()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")
}

func Test_Token_Create_Fail(t *testing.T) {

  x := NewTestToken()
  x.S3Bucket = ""
  err := x.Save()
  refute(t, err, nil)
  expect(t, x.Id, "")
}

func Test_Token_Update_Success(t *testing.T) {

  x := NewTestToken()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Save()
  expect(t, err, nil)
}

func Test_Token_Update_Fail(t *testing.T) {

  x := NewTestToken()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  x.S3Bucket = ""
  err = x.Save()
  refute(t, err, nil)
}

func Test_Token_Delete(t *testing.T) {

  x := NewTestToken()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Delete()
  expect(t, err, nil)
}

///////////

func Test_Token_BeforeCreate(t *testing.T) {
  x := NewTestToken()
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
}

func Test_Token_BeforeUpdate(t *testing.T) {
  x := NewTestToken()
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Token_UpdateAttrs(t *testing.T) {
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
    S3AccessId: attrs.S3AccessId,
    S3ClientSecretCrypted: u.Encrypt(tokenKey, attrs.S3ClientSecret),
    S3Bucket: attrs.S3Bucket,
    S3Region: attrs.S3Region,
  }

  expect(t, reflect.DeepEqual(targetByHand, obj), true)
}

func Test_TokenAttrs_Token(t *testing.T) {
  obj := &TokenAttrs{
    S3AccessId: "cheese",
    S3ClientSecret: "cheese",
    S3Bucket: "cheese",
    S3Region: "cheese",
  }
  targetByMethod := obj.Token()
  targetByHand := &Token{
    S3AccessId: obj.S3AccessId,
    S3ClientSecretCrypted: u.Encrypt(tokenKey, obj.S3ClientSecret),
    S3Bucket: obj.S3Bucket,
    S3Region: obj.S3Region,
  }

  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}



