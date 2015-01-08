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
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["JobKey"], true)

  x.JobKey = "cheese"
  x.Save()
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["JobKey"], false)
}

//////////////////////////////
// TRANSACTIONS //////////////

func Test_Item_Create_Success(t *testing.T) {

  x := NewTestItem()
  err := x.Save()
  expect(t, err, nil)
}

func Test_Item_Create_Fail(t *testing.T) {

  x := NewTestItem()
  x.JobKey = ""
  err := x.Save()
  refute(t, err, nil)
}

func Test_Item_Delete(t *testing.T) {

  x := NewTestItem()
  err := x.Save()
  expect(t, err, nil)

  err = x.Delete()
  expect(t, err, nil)
}

///////////

func Test_Item_BeforeCreate(t *testing.T) {
  x := NewTestItem()
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
}

