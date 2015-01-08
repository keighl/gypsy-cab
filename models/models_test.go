package models

import (
  "testing"
  "reflect"
  "image"
  u "gypsy/utils"
  r "github.com/dancannon/gorethink"
)

var (
  imgSquare *image.RGBA = image.NewRGBA(image.Rect(0, 0, 640, 640))
  imgLand *image.RGBA = image.NewRGBA(image.Rect(0, 0, 640, 480))
  imgPortrait *image.RGBA = image.NewRGBA(image.Rect(0, 0, 480, 500))
)

func init() {
  Config = u.ConfigForFile("../config/test.json")
  DB = u.RethinkSession(Config)

  _, _ = r.Table("users").Delete().RunWrite(DB)
  _, _ = r.Table("password_resets").Delete().RunWrite(DB)
  _, _ = r.Table("tokens").Delete().RunWrite(DB)
  _, _ = r.Table("jobs").Delete().RunWrite(DB)
  _, _ = r.Table("items").Delete().RunWrite(DB)
}

func expect(t *testing.T, a interface{}, b interface{}) {
  if a != b {
    t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func refute(t *testing.T, a interface{}, b interface{}) {
  if a == b {
    t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}
