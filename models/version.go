package models

import (
  "image"
  "image/draw"
  "regexp"
  "strconv"
  "github.com/disintegration/gift"
)

type Version struct {
  Key string `gorethink:"key" json:"key"`
  Recipe Recipe `gorethink:"recipe" json:"recipe"`
  Path string `gorethink:"path" json:"path"`
}

func (x *Version) MimeType() string {
  var mimeType string
  switch x.Recipe.Format {
    case "png":
      mimeType = "image/png"
    case "gif":
      mimeType = "image/gif"
    default:
      mimeType = "image/jpeg"
  }
  return mimeType
}

/////////////////

func (v *Version) ProcessRecipe(oImg image.Image) (img image.Image, err error) {
  img = oImg

  recipesX, _ := regexp.Compile(`(\w+\([a-z0-9,]*\))`)
  cropX, _ := regexp.Compile(`crop\((\d+),(\d+)\)`)
  resizeX, _ := regexp.Compile(`resize\((\d+),(\d+)\)`)
  grayX, _ := regexp.Compile(`grayscale\(\)`)
  blurX, _ := regexp.Compile(`blur\(([a-z0-9\.]+)\)`)

  recipesR := recipesX.FindAllStringSubmatch(v.Recipe.Instructions, -1)

  for _ , recipeMatch := range recipesR {
    hit := false

    if (!hit) {
      resizeR := resizeX.FindStringSubmatch(recipeMatch[0])
      if (len(resizeR) == 3) {
        // [resize(50,50) 50 50]
        width, _ := strconv.ParseUint(resizeR[1], 0, 64)
        height, _ := strconv.ParseUint(resizeR[2], 0, 64)
        img = Resize(img, int(width), int(height))
        hit = true
      }
    }

    if (!hit) {
      cropR := cropX.FindStringSubmatch(recipeMatch[0])
      if (len(cropR) == 3) {
        // [crop(50,50) 50 50]
        width, _ := strconv.ParseUint(cropR[1], 0, 64)
        height, _ := strconv.ParseUint(cropR[2], 0, 64)
        img = Crop(img, int(width), int(height))
        hit = true
      }
    }

    if (!hit) {
      grayR := grayX.FindStringSubmatch(recipeMatch[0])
      if (len(grayR) > 0) {
        // [grayscale()]
        img = Grayscale(img)
        hit = true
      }
    }

    if (!hit) {
      blurR := blurX.FindStringSubmatch(recipeMatch[0])
      if (len(blurR) == 2) {
        // [blur(2) 2]
        sigma, _ := strconv.ParseFloat(blurR[1], 64)
        img = Blur(img, float32(sigma))
        hit = true
      }
    }
  }

  return img, err
}

func Resize(imgIn image.Image, width int, height int) (image.Image) {
  var drawOut draw.Image
  filter := gift.Resize(width, height, gift.LanczosResampling)
  drawOut = image.NewRGBA(filter.Bounds(imgIn.Bounds()))
  filter.Draw(drawOut, imgIn, nil)
  return drawOut
}

func Grayscale(imgIn image.Image) (image.Image) {
  var drawOut draw.Image
  filter := gift.Grayscale()
  drawOut = image.NewRGBA(filter.Bounds(imgIn.Bounds()))
  filter.Draw(drawOut, imgIn, nil)
  return drawOut
}

func Blur(imgIn image.Image, sigma float32) (image.Image) {
  var drawOut draw.Image
  filter := gift.GaussianBlur(sigma)
  drawOut = image.NewRGBA(filter.Bounds(imgIn.Bounds()))
  filter.Draw(drawOut, imgIn, nil)
  return drawOut
}

func Crop(imgIn image.Image, width int, height int) (image.Image) {

  inW := imgIn.Bounds().Dx()
  inH := imgIn.Bounds().Dy()

  inAspect := float64(inW) / float64(inH)
  outAspect := float64(width) / float64(height)

  if inAspect > outAspect {
    imgIn = Resize(imgIn, 0, height)
  } else {
    imgIn = Resize(imgIn, width, 0)
  }

  inW = imgIn.Bounds().Dx()
  inH = imgIn.Bounds().Dy()

  inMinX := imgIn.Bounds().Min.X
  inMinY := imgIn.Bounds().Min.Y

  centerX := inMinX + inW/2
  centerY := inMinY + inH/2

  x0 := centerX - width/2
  y0 := centerY - height/2
  x1 := x0 + width
  y1 := y0 + height

  var drawOut draw.Image
  filter := gift.Crop(image.Rect(x0, y0, x1, y1))
  drawOut = image.NewRGBA(filter.Bounds(imgIn.Bounds()))
  filter.Draw(drawOut, imgIn, nil)
  return drawOut
}


