package models

import (
  "time"
  r "github.com/dancannon/gorethink"
  "errors"
  "regexp"
  "strings"
  "github.com/dchest/uniuri"
)

type Job struct {
  Errors []string `gorethink:"-" json:"errors,omitempty"`
  ErrorMap map[string]bool `gorethink:"-" json:"-"`
  Id string `gorethink:"id,omitempty" json:"-"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at,omitempty"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"updated_at,omitempty"`
  UserId string `gorethink:"user_id" json:"-"`
  Key string `gorethink:"key" json:"key"`
  Recipes []Recipe `gorethink:"recipes" json:"recipes"`
}

type JobAttrs struct {
  Key string `json:"key" form:"key"`
  Recipes []Recipe `json:"recipes" form:"recipes"`
}

func (job *Job) BuildItem() *Item {
  item := &Item{
    JobKey: job.Key,
    JobId: job.Id,
    Id: uniuri.NewLen(15),
    JobUserId: job.UserId,
  }

  for _, r := range job.Recipes {
    v := Version{
      Key: r.Key,
      Recipe: r,
      Path: job.Key + "/" + item.Id + "/" + r.Key + r.Extension(),
    }
    item.Versions = append(item.Versions, v)
  }
  return item
}

//////////////////////////////
// TRANSACTIONS //////////////

func (x *Job) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  if (x.Id == "") {
    x.BeforeCreate()
    res, err := r.Table("jobs").Insert(x).RunWrite(DB)
    if (err != nil) { return err }
    x.Id = res.GeneratedKeys[0]
  }

  x.BeforeUpdate()
  _, err := r.Table("jobs").Get(x.Id).Replace(x).RunWrite(DB)
  return err
}

func (x *Job) Delete() error {
  _, err := r.Table("jobs").Get(x.Id).Delete().RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *Job) BeforeCreate() {
  x.CreatedAt = time.Now()
  x.UpdatedAt = time.Now()
}

func (x *Job) BeforeUpdate() {
  x.UpdatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Job) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.Trimspace()
  x.ValidateRecipes()
  x.ValidateKey()
  return !x.HasErrors()
}

func (x *Job) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Job) Trimspace() {
  x.Key = strings.TrimSpace(x.Key)
  x.Key = strings.ToLower(x.Key)
}

func (x *Job) ValidateRecipes() {
  if (len(x.Recipes) == 0) {
    x.Errors = append(x.Errors, "You must assign some recipes")
    x.ErrorMap["Recipes"] = true
  }

  for _, r := range x.Recipes {
    if (!r.Validate()) {
      x.Errors = append(x.Errors, r.Errors...)
      x.ErrorMap["Recipes"] = true
    }
  }
}

func (x *Job) ValidateKey() {
  regex := regexp.MustCompile(`^[a-z0-9\-_]{4,}$`)
  if (!regex.MatchString(x.Key)) {
    x.Errors = append(x.Errors, "Key must be at least 4 characters, and contain only letters, numbers, dashes or underscores")
    x.ErrorMap["Key"] = true
  }
}

//////////////////////////////
// OTHER /////////////////////

func (x *Job) UpdateFromAttrs(attrs JobAttrs) {
  // Only recipes... key isn't mutable
  if (len(attrs.Recipes) > 0) { x.Recipes = attrs.Recipes }
}

func (x *JobAttrs) Job() (*Job) {
  return &Job{
    Key: x.Key,
    Recipes: x.Recipes,
  }
}
