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

func Test_Job_Table(t *testing.T) {
  x := NewTestJob()
  expect(t, x.Table(), "jobs")
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
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Recipes"], true)

  x = NewTestJob()
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Recipes"], false)
}

func Test_Job_RequiresValidRecipes(t *testing.T) {

  x := NewTestJob()
  x.Recipes = []Recipe{Recipe{Key: ""}}
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Recipes"], true)

  x = NewTestJob()
  x.Recipes = []Recipe{NewTestRecipe()}
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Recipes"], false)
}

func Test_Job_RequiresKey(t *testing.T) {
  x := NewTestJob()
  x.Key = ""
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Key"], true)
}

func Test_Job_RequiresKeyForm(t *testing.T) {
  x := NewTestJob()
  x.Key = "cheese asdflkj"
  err := Save(x)
  refute(t, err, nil)
  expect(t, x.ErrorMap["Key"], true)

  x.Key = "cheese-thing_ghi"
  err = Save(x)
  expect(t, err, nil)
  expect(t, x.ErrorMap["Key"], false)
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

