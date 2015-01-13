package models

import (
  "testing"
)

func NewTestItem() *Item {
  j := NewTestJobPersisted()
  x := j.BuildItem()
  return x
}

//////////////////////////////
// VALIDATIONS ///////////////

func Test_Item_RequiresJobKey(t *testing.T) {

  x := NewTestItem()

  x.JobKey = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["JobKey"], true)

  x.JobKey = "cheese"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["JobKey"], false)
}

