package api

import (
  "testing"
  m "github.com/keighl/gypsy-cab/models"
  "net/http"
  "errors"
)

//////////////////////////////////////
// INDEX ///////////////////

func itemsIndexRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/jobs/:job_id/items", AuthorizeOK, ItemsIndex)
  req, _ := http.NewRequest("GET", "/v1/jobs/XXXXX/items", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Items_Index_JobError(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXX")
    return nil, errors.New("******")
  }
  itemsIndexRunner(t, http.StatusNotFound)
}

func Test_Route_Items_Index_Error(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXX")
    return &m.Job{}, nil
  }
  loadItems = func(u *m.Job) ([]m.Item, error) {
    return nil, errors.New("*****")
  }
  itemsIndexRunner(t, http.StatusInternalServerError)
}

func Test_Route_Items_Index_Success(t *testing.T) {
  loadJob = func(id string, u *m.User) (*m.Job, error) {
    expect(t, id, "XXXXX")
    return &m.Job{}, nil
  }
  loadItems = func(u *m.Job) ([]m.Item, error) {
    return []m.Item{m.Item{}}, nil
  }
  itemsIndexRunner(t, http.StatusOK)
}

//////////////////////////////////////
// SHOW ///////////////////

func itemsShowRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/items/:item_id", AuthorizeOK, ItemsShow)
  req, _ := http.NewRequest("GET", "/v1/items/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Items_Show_NotFound(t *testing.T) {
  loadItem = func(id string, u *m.User) (*m.Item, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }
  itemsShowRunner(t, http.StatusNotFound)
}

func Test_Route_Items_Show_Success(t *testing.T) {
  loadItem = func(id string, u *m.User) (*m.Item, error) {
    expect(t, id, "XXXXXX")
    return &m.Item{}, nil
  }

  itemsShowRunner(t, http.StatusOK)
}

//////////////////////////////////////
// DELETE ///////////////////

func itemsDeleteRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Delete("/v1/items/:item_id", AuthorizeOK, ItemsDelete)
  req, _ := http.NewRequest("DELETE", "/v1/items/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Items_Delete_NotFound(t *testing.T) {
  loadItem = func(id string, u *m.User) (*m.Item, error) {
    expect(t, id, "XXXXXX")
    return nil, nil
  }

  itemsDeleteRunner(t, http.StatusNotFound)
}

func Test_Route_Items_Delete_Success(t *testing.T) {
  loadItem = func(id string, u *m.User) (*m.Item, error) {
    expect(t, id, "XXXXXX")
    return &m.Item{}, nil
  }

  deleteItem = func(item *m.Item) (error) {
    return nil
  }

  itemsDeleteRunner(t, http.StatusOK)
}

func Test_Route_Items_Delete_400(t *testing.T) {
  loadItem = func(id string, u *m.User) (*m.Item, error) {
    expect(t, id, "XXXXXX")
    return &m.Item{}, nil
  }

  deleteItem = func(item *m.Item) (error) {
    return errors.New("*******")
  }

  itemsDeleteRunner(t, http.StatusBadRequest)
}

