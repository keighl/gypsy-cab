package models

import (
  "time"
  r "github.com/dancannon/gorethink"
  "errors"
)

//////////////////////////////
// ITEM //////////////////////

type Item struct {
  Errors []string `gorethink:"-" json:"errors,omitempty" sql:"-"`
  ErrorMap map[string]bool `gorethink:"-" json:"-" sql:"-"`
  Id string `gorethink:"id" json:"id"`
  CreatedAt time.Time `gorethink:"created_at" json:"created_at,omitempty"`
  UpdatedAt time.Time `gorethink:"updated_at" json:"updated_at,omitempty"`
  TokenId string `gorethink:"token_id" json:"-"`
  JobId string `gorethink:"job_id" json:"-"`
  JobKey string `gorethink:"job_key" json:"job"`
  JobUserId string `gorethink:"job_user_id" json:"-"`
  Versions []Version `gorethink:"versions" json:"versions"`
}

//////////////////////////////
// TRANSACTIONS //////////////

// No updates
func (x *Item) Save() error {

  if (!x.Validate()) {
    return errors.New("Validation errors")
  }

  x.BeforeCreate()
  _, err := r.Table("items").Insert(x).RunWrite(DB)
  if (err != nil) { return err }

  return nil
}

func (x *Item) Delete() error {
  _, err := r.Table("items").Get(x.Id).Delete().RunWrite(DB)
  return err
}

//////////////////////////////
// CALLBACKS /////////////////

func (x *Item) BeforeCreate() {
  x.CreatedAt = time.Now()
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Item) Validate() (bool) {
  x.Errors = []string{}
  x.ErrorMap = map[string]bool{}
  x.ValidateJobKey()
  return !x.HasErrors()
}

func (x *Item) HasErrors() (bool) {
  return len(x.Errors) > 0
}

func (x *Item) ValidateJobKey() {

  if (x.JobKey == "") {
    x.Errors = append(x.Errors, "You must specify the Job key")
    x.ErrorMap["JobKey"] = true
  }
}


