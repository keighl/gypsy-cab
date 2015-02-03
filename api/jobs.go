package api

import (
  m "github.com/keighl/gypsy-cab/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "errors"
)

type JobsData struct {
  User *m.User `json:"current_user,omitempty"`
  Jobs []m.Job `json:"jobs"`
}

type JobData struct {
  User *m.User `json:"current_user,omitempty"`
  *m.Job `json:"job"`
}

////////////////////////

var loadJobs = func(u *m.User) ([]m.Job, error) {
  jobs := []m.Job{}
  res, err := r.Table("jobs").
    GetAllByIndex("user_id", u.Id).
    OrderBy(r.Desc("created_at")).
    Run(DB)
  if (err != nil) { return nil, err }
  err = res.All(&jobs)
  return jobs, err
}

func JobsIndex(r render.Render, user *m.User) {
  jobs, err := loadJobs(user)

  if (err != nil) {
    r.JSON(500, MessageEnvelope(err.Error()))
    return
  }

  data := &JobsData{User: user, Jobs: jobs}
  r.JSON(200, data)
}

/////////////////////////

var loadJob = func(id string, u *m.User) (*m.Job, error) {
  job := &m.Job{}
  res, err := r.Table("jobs").GetAllByIndex("key", id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(job)
  if (err != nil) { return nil, err }
  if (job.UserId != u.Id) { return nil, errors.New("Not your job") }
  return job, err
}

func JobsShow(params martini.Params, r render.Render, user *m.User) {
  job, err := loadJob(params["job_id"], user)

  if (job == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  data := &JobData{User: user, Job: job}
  r.JSON(200, data)
}

/////////////////////////

var saveJob = func(job *m.Job) (error) {
  return m.Save(job)
}

func JobsCreate(r render.Render, user *m.User, attrs m.JobAttrs) {
  job := attrs.Job()
  job.UserId = user.Id
  err := saveJob(job)

  if (err != nil) {
    if (job.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), job.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &JobData{User: user, Job: job}
  r.JSON(201, data)
}

/////////////////////////

func JobsUpdate(params martini.Params, r render.Render, user *m.User, attrs m.JobAttrs) {
  job, err := loadJob(params["job_id"], user)

  if (job == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  job.UpdateFromAttrs(attrs)
  err = saveJob(job)

  if (err != nil) {
    if (job.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), job.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &JobData{User: user, Job: job}
  r.JSON(200, data)
}

/////////////////////////

var deleteJob = func(job *m.Job) (error) {
  return m.Delete(job)
}

func JobsDelete(params martini.Params, r render.Render, user *m.User) {
  job, err := loadJob(params["job_id"], user)

  if (job == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  err = deleteJob(job)

  if (err != nil) {
    r.JSON(400, MessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &Data{User: user, Message: &Message{"The job was deleted"}}
  r.JSON(200, data)
}
