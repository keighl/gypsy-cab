package api

import (
  m "github.com/keighl/gypsy-cab/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "errors"
)

type ItemsData struct {
  User *m.User `json:"current_user,omitempty"`
  *m.Job `json:"job,omitempty"`
  Items []m.Item `json:"item,omitempty"`
}

type ItemData struct {
  User *m.User `json:"current_user,omitempty"`
  *m.Item `json:"item,omitempty"`
}

////////////////////////

var loadItems = func(job *m.Job) ([]m.Item, error) {
  items := []m.Item{}
  res, err := r.Table("items").
    GetAllByIndex("job_id", job.Id).
    OrderBy(r.Desc("created_at")).
    Run(DB)
  if (err != nil) { return nil, err }
  err = res.All(&items)
  return items, err
}

func ItemsIndex(params martini.Params, r render.Render, user *m.User) {

  job, err := loadJob(params["job_id"], user)

  if (err != nil) {
    r.JSON(404, MessageEnvelope("Job not found"))
    return
  }

  items, err := loadItems(job)

  if (err != nil) {
    r.JSON(500, MessageEnvelope(err.Error()))
    return
  }

  data := &ItemsData{User: user, Items: items, Job: job}
  r.JSON(200, data)
}

/////////////////////////

var loadItem = func(id string, u *m.User) (*m.Item, error) {
  item := &m.Item{}
  res, err := r.Table("items").Get(id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(item)
  if (err != nil) { return nil, err }
  if (item.JobUserId != u.Id) { return nil, errors.New("Not your item") }
  return item, err
}

func ItemsShow(params martini.Params, r render.Render, user *m.User) {
  item, err := loadItem(params["item_id"], user)

  if (item == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  data := &ItemData{User: user, Item: item}
  r.JSON(200, data)
}

/////////////////////////

var deleteItem = func(item *m.Item) (error) {
  return m.Delete(item)
}

func ItemsDelete(params martini.Params, r render.Render, user *m.User) {
  item, err := loadItem(params["item_id"], user)

  if (item == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  err = deleteItem(item)

  if (err != nil) {
    r.JSON(400, MessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &Data{User: user, Message: &Message{"The item was deleted"}}
  r.JSON(200, data)
}
