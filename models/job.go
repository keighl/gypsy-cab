package models

import (
  "strings"
  "github.com/dchest/uniuri"
  "regexp"
)

type Job struct {
  Record
  Id string `gorethink:"id,omitempty" json:"id"`
  UserId string `gorethink:"user_id" json:"-"`
  Key string `gorethink:"key" json:"key"`
  Recipes []Recipe `gorethink:"recipes" json:"recipes"`
}

type JobAttrs struct {
  Key string `json:"key" form:"key"`
  Recipes []Recipe `json:"recipes" form:"recipes"`
}

func (x *Job) Table() string {
  return "jobs"
}

func (job *Job) BuildItem() *Item {
  item := &Item{
    Id: uniuri.NewLen(15),
    JobKey: job.Key,
    JobId: job.Id,
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
// VALIDATIONS ///////////////

func (x *Job) Validate() {
  x.Record.Validate()
  x.Trimspace()
  x.ValidateRecipes()
  x.ValidateKey()
}

func (x *Job) ValidateRecipes() {
  if (len(x.Recipes) == 0) {
    x.ErrorOn("Recipes", "You must assign some recipes")
  }

  for _, r := range x.Recipes {
    r.Validate()
    if (r.HasErrors()) {
      x.ErrorOn("Recipes", r.Errors...)
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

func (x *Job) Trimspace() {
  x.Key = strings.TrimSpace(x.Key)
  x.Key = strings.ToLower(x.Key)
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
