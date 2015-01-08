package models

import (
  "testing"
)

func NewTestVersion() Version{
  return Version{
    Recipe: NewTestRecipe(),
  }
}

func Test_Version_MimeType(t *testing.T) {
  x := NewTestVersion()

  x.Recipe.Format = "png"
  expect(t, x.MimeType(), "image/png")

  x.Recipe.Format = "gif"
  expect(t, x.MimeType(), "image/gif")

  x.Recipe.Format = "jpg"
  expect(t, x.MimeType(), "image/jpeg")

  x.Recipe.Format = "jpeg"
  expect(t, x.MimeType(), "image/jpeg")

  x.Recipe.Format = "cheese"
  expect(t, x.MimeType(), "image/jpeg")
}

func Test_Version_Process(t *testing.T) {
  recipes := []string{"", "resize(50,50)", "crop(50,50)", "grayscale()", "blur(2)"}

  for _, r := range recipes {
    v := Version{
      Key: "XXX",
      Recipe: Recipe{Format: "png", Instructions: r, Key: "XXX"},
      Path: "XXXXXX/XXXXXX/XXX.png",
    }

    imgRes, err := v.ProcessRecipe(imgSquare)
    refute(t, &imgRes, nil)
    expect(t, err, nil)
  }
}

func Test_Version_CropLand(t *testing.T) {
  v := Version{
    Key: "XXX",
    Recipe: Recipe{Format: "png", Instructions: "crop(50,50)", Key: "XXX"},
    Path: "XXXXXX/XXXXXX/XXX.png",
  }

  imgRes, err := v.ProcessRecipe(imgLand)
  refute(t, &imgRes, nil)
  expect(t, err, nil)
}