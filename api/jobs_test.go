package api

import (
  "testing"
  m "gypsy/models"
  "net/http"
  "errors"
  "bytes"
  "encoding/json"
  "github.com/martini-contrib/binding"
)

//////////////////////////////////////
// INDEX ///////////////////

func jobsIndexRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/jobs", AuthorizeOK, JobsIndex)
  req, _ := http.NewRequest("GET", "/v1/jobs", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Jobs_Index_Error(t *testing.T) {
  loadJobs = func(u *m.User) ([]m.Job, error) {
    return nil, errors.New("*****")
  }
  jobsIndexRunner(t, http.StatusInternalServerError)
}

func Test_Route_Jobs_Index_Success(t *testing.T) {
  loadJobs = func(u *m.User) ([]m.Job, error) {
    return []m.Job{m.Job{}}, nil
  }
  jobsIndexRunner(t, http.StatusOK)
}

//////////////////////////////////////
// SHOW ///////////////////

func jobsShowRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/jobs/:job_id", AuthorizeOK, JobsShow)
  req, _ := http.NewRequest("GET", "/v1/jobs/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Jobs_Show_NotFound(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }
  jobsShowRunner(t, http.StatusNotFound)
}

func Test_Route_Jobs_Show_Success(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  jobsShowRunner(t, http.StatusOK)
}

//////////////////////////////////////
// CREATE ///////////////////

func jobsCreateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Post("/v1/jobs", AuthorizeOK, binding.Bind(m.JobAttrs{}), JobsCreate)
  body, _ := json.Marshal(m.JobAttrs{})
  req, _ := http.NewRequest("POST", "/v1/jobs", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Jobs_Create_Failure_400(t *testing.T) {
  saveJob = func(job *m.Job) (error) {
    job.Errors = []string{"Something went wrong!"}
    return errors.New("*******")
  }
  jobsCreateRunner(t, http.StatusBadRequest)
}

func Test_Route_Jobs_Create_Failure_500(t *testing.T) {
  saveJob = func(job *m.Job) (error) {
    return errors.New("*******")
  }
  jobsCreateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Jobs_Create_Success(t *testing.T) {
  saveJob = func(job *m.Job) (error) {
    return nil
  }
  jobsCreateRunner(t, http.StatusCreated)
}

//////////////////////////////////////
// UPDATE ///////////////////

func jobsUpdateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Put("/v1/jobs/:job_id", AuthorizeOK, binding.Bind(m.JobAttrs{}), JobsUpdate)
  body, _ := json.Marshal(m.JobAttrs{})
  req, _ := http.NewRequest("PUT", "/v1/jobs/XXXXXX", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Jobs_Update_NotFound(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }

  jobsUpdateRunner(t, http.StatusNotFound)
}

func Test_Route_Jobs_Update_Failure_400(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  saveJob = func(job *m.Job) (error) {
    job.Errors = []string{"Something went wrong!"}
    return errors.New("*********")
  }

  jobsUpdateRunner(t, http.StatusBadRequest)
}

func Test_Route_Jobs_Update_Failure_500(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  saveJob = func(job *m.Job) (error) {
    return errors.New("*********")
  }

  jobsUpdateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Jobs_Update_Success(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  saveJob = func(job *m.Job) (error) {
    return nil
  }

  jobsUpdateRunner(t, http.StatusOK)
}

//////////////////////////////////////
// DELETE ///////////////////

func jobsDeleteRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Delete("/v1/jobs/:job_id", AuthorizeOK, JobsDelete)
  req, _ := http.NewRequest("DELETE", "/v1/jobs/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Jobs_Delete_NotFound(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return nil, nil
  }

  jobsDeleteRunner(t, http.StatusNotFound)
}

func Test_Route_Jobs_Delete_Success(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  deleteJob = func(job *m.Job) (error) {
    return nil
  }

  jobsDeleteRunner(t, http.StatusOK)
}

func Test_Route_Jobs_Delete_400(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXXX")
    return &m.Job{}, nil
  }

  deleteJob = func(job *m.Job) (error) {
    return errors.New("*******")
  }

  jobsDeleteRunner(t, http.StatusBadRequest)
}

