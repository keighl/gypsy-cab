package models

import (
  "strings"
  "regexp"
)

type Recipe struct {
  Errors []string `gorethink:"-" json:"errors,omitempty"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
  Key string `gorethink:"key" json:"key" form:"key"`
  Instructions string `gorethink:"instructions" json:"instructions" form:"instructions"`
  Format string `gorethink:"format" json:"format" form:"format"`
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Recipe) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.Trimspace()
  x.ValidateKey()
  x.ValidateFormat()
  return !x.HasErrors()
}

func (x *Recipe) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Recipe) Trimspace() {
  x.Key          = strings.TrimSpace(x.Key)
  x.Instructions = strings.TrimSpace(x.Instructions)
  x.Format       = strings.TrimSpace(x.Format)

  x.Key          = strings.ToLower(x.Key)
  x.Instructions = strings.ToLower(x.Instructions)
  x.Format       = strings.ToLower(x.Format)
}

func (x *Recipe) ValidateKey() {
  regex := regexp.MustCompile(`^[a-z0-9\-_]{4,}$`)
  if (!regex.MatchString(x.Key)) {
    x.Errors = append(x.Errors, "Key must be at least 4 characters, and contain only letters, numbers, dashes or underscores")
    x.ErrorMap["Key"] = true
  }
}

func (x *Recipe) ValidateFormat() {

  if (x.Format != "") {
    if (!stringInSlice(x.Format, []string{"jpg", "jpeg", "gif", "png"})) {
      x.Errors = append(x.Errors, "Format must be jpg, gif, or png")
      x.ErrorMap["Format"] = true
    }
  }
}

func stringInSlice(a string, list []string) bool {
  for _, b := range list {
    if b == a {
      return true
    }
  }
  return false
}

func (x *Recipe) Extension() (ext string) {
  switch x.Format {
    case "png":
      ext = ".png"
    case "gif":
      ext = ".gif"
    default:
      ext = ".jpg"
  }
  return ext
}
