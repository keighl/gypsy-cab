package models


//////////////////////////////
// ITEM //////////////////////

type Item struct {
  Record
  Id string `gorethink:"id,omitempty" json:"id"`
  TokenId string `gorethink:"token_id" json:"-"`
  JobId string `gorethink:"job_id" json:"-"`
  JobKey string `gorethink:"job_key" json:"job"`
  JobUserId string `gorethink:"job_user_id" json:"-"`
  Versions []Version `gorethink:"versions" json:"versions"`
}

func (x *Item) IsNewRecord() bool {
  // Items are always new! i.e. no updating
  return true
}

//////////////////////////////
// VALIDATIONS ///////////////

func (x *Item) Validate() {
  x.Record.Validate()
  x.ValidateJobKey()
}

func (x *Item) ValidateJobKey() {
  if (x.JobKey == "") {
    x.ErrorOn("JobKey", "You must specify the Job key")
  }
}


