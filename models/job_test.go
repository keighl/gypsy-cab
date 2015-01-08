package models

import (
  "testing"
  "reflect"
  "github.com/dchest/uniuri"
)

func NewTestJob() *Job {
  return &Job{
    Key: uniuri.NewLen(15),
    Recipes: []Recipe{NewTestRecipe(), NewTestRecipe()},
  }
}

func NewTestJobPersisted() *Job {
  x := NewTestJob()
  x.Id = uniuri.NewLen(30)
  return x
}

func Test_Item_BuildItem(t *testing.T) {

  j := NewTestJobPersisted()
  x := j.BuildItem()

  expect(t, x.JobKey, j.Key)
  expect(t, x.JobId, j.Id)
  expect(t, len(x.Versions), len(j.Recipes))

  for i, r := range j.Recipes {
    expect(t, x.Versions[i].Path, j.Key + "/" + x.Id + "/" + r.Key + r.Extension())
    expect(t, x.Versions[i].Key, r.Key)
    expect(t, reflect.DeepEqual(x.Versions[i].Recipe, r), true)
  }
}

//////////////////////////////
// VALIDATIONS ///////////////

func Test_Job_RequiresRecipes(t *testing.T) {

  x := NewTestJob()
  x.Recipes = []Recipe{}
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Recipes"], true)

  x = NewTestJob()
  x.Save()
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Recipes"], false)
}

func Test_Job_RequiresValidRecipes(t *testing.T) {

  x := NewTestJob()
  x.Recipes = []Recipe{Recipe{Key: ""}}
  x.Save()
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Recipes"], true)

  x = NewTestJob()
  x.Recipes = []Recipe{NewTestRecipe()}
  x.Save()
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Recipes"], false)
}

func Test_Job_RequiresKey(t *testing.T) {
  x := NewTestJob()
  x.Key = ""
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Key"], true)
}

func Test_Job_RequiresKeyForm(t *testing.T) {
  x := NewTestJob()
  x.Key = "cheese asdflkj"
  expect(t, x.Validate(), false)
  expect(t, x.ErrorMap["Key"], true)

  x.Key = "cheese-thing_ghi"
  expect(t, x.Validate(), true)
  expect(t, x.ErrorMap["Key"], false)
}

//////////////////////////////
// TRANSACTIONS //////////////

func Test_Job_Create_Success(t *testing.T) {

  x := NewTestJob()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")
}

func Test_Job_Create_Fail(t *testing.T) {

  x := NewTestJob()
  x.Recipes = []Recipe{}
  err := x.Save()
  refute(t, err, nil)
  expect(t, x.Id, "")
}

func Test_Job_Update_Success(t *testing.T) {

  x := NewTestJobPersisted()
  err := x.Save()
  expect(t, err, nil)
}

func Test_Job_Update_Fail(t *testing.T) {

  x := NewTestJobPersisted()
  x.Recipes = []Recipe{}
  err := x.Save()
  refute(t, err, nil)
}

func Test_Job_Delete(t *testing.T) {

  x := NewTestJob()
  err := x.Save()
  expect(t, err, nil)
  refute(t, x.Id, "")

  err = x.Delete()
  expect(t, err, nil)
}

///////////

func Test_Job_BeforeCreate(t *testing.T) {
  x := NewTestJob()
  x.BeforeCreate()
  refute(t, x.CreatedAt.Format("RFC3339"), nil)
}

func Test_Job_BeforeUpdate(t *testing.T) {
  x := NewTestJob()
  x.BeforeUpdate()
  refute(t, x.UpdatedAt.Format("RFC3339"), nil)
}

/////////////////////////
// ATTR CONVERSION //////

func Test_Job_UpdateAttrs(t *testing.T) {
  obj := &Job{
    Recipes: []Recipe{NewTestRecipe()},
  }
  r := NewTestRecipe()
  r.Key = "keyyyyy"
  attrs := JobAttrs{
    Recipes: []Recipe{r},
  }
  obj.UpdateFromAttrs(attrs)
  targetByHand := &Job{
    Recipes: attrs.Recipes,
  }

  expect(t, reflect.DeepEqual(targetByHand, obj), true)
}

func Test_JobAttrs_Job(t *testing.T) {
  obj := &JobAttrs{}
  targetByMethod := obj.Job()
  targetByHand := &Job{}
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}

