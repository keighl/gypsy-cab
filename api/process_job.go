package api

import (
  m "gypsy/models"
  r "github.com/dancannon/gorethink"
  "os"
  "fmt"
  "bufio"
  "github.com/martini-contrib/render"
  "github.com/go-martini/martini"
  "mime/multipart"
  "github.com/disintegration/gift"
  "github.com/rwcarlsen/goexif/exif"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "image"
  "image/draw"
  "image/jpeg"
  "image/png"
  "image/gif"
)

type Upload struct {
  Image *multipart.FileHeader `form:"image"`
}

var loadProcessToken = func(t string) (*m.Token, error) {
  token := &m.Token{}
  res, err := r.Table("tokens").Get(t).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(token)
  if (err != nil) { return nil, err }
  return token, err
}

var loadProcessJob = func(id string, t *m.Token) (*m.Job, error) {
  job := &m.Job{}
  filter := func(row r.Term) r.Term {
    return row.Field("user_id").Eq(t.UserId).And(row.Field("key").Eq(id))
  }
  res, err := r.Table("jobs").
    Filter(filter).
    Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(&job)
  return job, err
}

func JobProcess(params martini.Params, r render.Render, upload Upload) {

  // Load the token
  token, err := loadProcessToken(params["token"])
  if (token == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Token not found"))
    return
  }

  // Load the job
  job, err := loadProcessJob(params["job_id"], token)
  if (job == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Job not found"))
    return
  }

  // Create an item
  item := job.BuildItem()

  // Open the uploaded file
  file, err := upload.Image.Open()
  if (err != nil) {
    r.JSON(400, MessageEnvelope("Unable to parse file"))
    return
  }
  defer file.Close()

  // Decode the file into an image
  originalImg, fileType, err := image.Decode(file)
  if (err != nil) {
    r.JSON(400, MessageEnvelope("Unable to parse file"))
    return
  }

  // Check for supported types
  switch fileType {
    case "png":
    case "jpeg":
    case "gif":
    default:
      if (err != nil) {
        r.JSON(400, MessageEnvelope("Must be JPG, PNG or GIF"))
        return
      }
  }

  // Fix EXIF rotation if applicable
  exifInfo, err := exif.Decode(file)
  if err != nil {
    // No exif data
  } else {
    orientation, err := exifInfo.Get(exif.Orientation)
    if err != nil {
      r.JSON(500, ServerErrorEnvelope())
      return
    }
    orientationInt, _ := orientation.Int(0)
    originalImg = FixExif(orientationInt, originalImg)
  }

  // For each version, process the image, and write it to a temp file
  for _, v := range item.Versions {
    img, err := v.ProcessRecipe(originalImg)
    if err != nil {
      r.JSON(500, ServerErrorEnvelope())
      return
    }

    tmpPath := "tmp/" + item.Id + "-" + v.Key

    tempImage, err := os.Create(tmpPath)
    if err != nil {
      r.JSON(500, ServerErrorEnvelope())
      return
    }

    switch v.Recipe.Format {
      case "png":
        err = png.Encode(tempImage, img)
      case "gif":
        err = gif.Encode(tempImage, img, &gif.Options{NumColors: 256})
      default:
        err = jpeg.Encode(tempImage, img, &jpeg.Options{Quality: 95})
    }

    tempImage.Close()
    if err != nil {
      r.JSON(500, ServerErrorEnvelope())
      return
    }
    img = nil
  }

  remProcessed := len(item.Versions)
  errChan := make(chan error)
  uploadedChan := make(chan m.Version)

  for _, v := range item.Versions {
    go SendImageToS3(item, v, token, uploadedChan, errChan)
  }

  for {
    select {
    case _ = <-uploadedChan:
      remProcessed--
    case _ = <-errChan:
      if err != nil {
        r.JSON(500, MessageEnvelope("Unable file to send to S3"))
        return
      }
      break
    }

    if remProcessed == 0 {
      // Save everything to the DB
      err = m.Save(item)
      if err != nil {
        r.JSON(500, ServerErrorEnvelope())
        return
      }
      break
    }
  }

  if (token.LegacyToken) {
    // For older clients that may expect a data.{} envelope
    item.Uuid = item.Id
    r.JSON(201, LegacyEnvelope{&Data{Item: item}})
  } else {
    r.JSON(201, &Data{Item: item})
  }
}

///////

func FixExif(orientation int, imgIn image.Image) (image.Image) {
  var filter gift.Filter
  var drawOut draw.Image

  switch orientation {
    case 8:
      fmt.Printf("%d ROTATE 90\n", orientation)
      filter = gift.Rotate90()
    case 3:
      fmt.Printf("%d ROTATE 180\n", orientation)
      filter = gift.Rotate180()
    case 6:
      fmt.Printf("%d ROTATE 270\n", orientation)
      filter = gift.Rotate270()
    default:
      fmt.Println("ROTATE 0")
      return imgIn
  }

  drawOut = image.NewRGBA(filter.Bounds(imgIn.Bounds()))
  filter.Draw(drawOut, imgIn, nil)
  return drawOut
}

///////

func SendImageToS3( item *m.Item, v m.Version, token *m.Token, uploadChan chan m.Version, errChar chan error) {

  tmpPath := "tmp/" + item.Id + "-" + v.Key
  file, err := os.Open(tmpPath)
  if err != nil {
    errChar <- err
    return
  }
  defer file.Close()

  buffy := bufio.NewReader(file)
  stat, err := file.Stat()
  if err != nil {
    errChar <- err
    return
  }

  size := stat.Size()

  var auth aws.Auth
  auth.AccessKey = token.S3AccessId
  auth.SecretKey = token.S3ClientSecret()

  s := s3.New(auth, aws.Regions[token.S3Region])
  bucket := s.Bucket(token.S3Bucket)

  err = bucket.PutReader(v.Path, buffy, size, v.MimeType(), s3.BucketOwnerFull)
  if err != nil {
    errChar <- err
    return
  }

  uploadChan <- v
}

