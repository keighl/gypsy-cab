package models

import (
  "testing"
  "github.com/dchest/uniuri"
)

func NewTestRecipe() Recipe{
  return Recipe{
    Key: uniuri.NewLen(10),
    Instructions: "crop(100,100)",
    Format: "gif",
  }
}

func Test_Recipe_RequiresKey(t *testing.T) {
  x := NewTestRecipe()
  x.Key = ""
  x.Validate()
  expect(t, x.ErrorMap["Key"], true)
}

func Test_Recipe_RequiresKeyForm(t *testing.T) {
  x := NewTestRecipe()
  x.Key = "cheese asdflkj"
  x.Validate()
  expect(t, x.ErrorMap["Key"], true)

  x.Key = "cheese-thing_ghi"
  x.Validate()
  expect(t, x.ErrorMap["Key"], false)
}

func Test_Recipe_RequiresFormatInclusion(t *testing.T) {
  x := NewTestRecipe()
  x.Format = "cheese"
  x.Validate()
  expect(t, x.ErrorMap["Format"], true)

  for _, q := range []string{"", "jpg", "jpeg", "gif", "png"} {
    x.Format = q
    x.Validate()
    expect(t, x.ErrorMap["Format"], false)
  }
}

func Test_Recipe_Extension(t *testing.T) {
  x := NewTestRecipe()

  x.Format = "png"
  expect(t, x.Extension(), ".png")

  x.Format = "gif"
  expect(t, x.Extension(), ".gif")

  x.Format = "jpg"
  expect(t, x.Extension(), ".jpg")

  x.Format = "jpeg"
  expect(t, x.Extension(), ".jpg")

  x.Format = "cheese"
  expect(t, x.Extension(), ".jpg")
}

