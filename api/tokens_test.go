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

func tokensIndexRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/tokens", AuthorizeOK, TokensIndex)
  req, _ := http.NewRequest("GET", "/v1/tokens", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Tokens_Index_Error(t *testing.T) {
  loadTokens = func(u *m.User) ([]m.Token, error) {
    return nil, errors.New("*****")
  }
  tokensIndexRunner(t, http.StatusInternalServerError)
}

func Test_Route_Tokens_Index_Success(t *testing.T) {
  loadTokens = func(u *m.User) ([]m.Token, error) {
    return []m.Token{m.Token{}}, nil
  }
  tokensIndexRunner(t, http.StatusOK)
}

//////////////////////////////////////
// SHOW ///////////////////

func tokensShowRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Get("/v1/tokens/:token", AuthorizeOK, TokensShow)
  req, _ := http.NewRequest("GET", "/v1/tokens/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Tokens_Show_NotFound(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }
  tokensShowRunner(t, http.StatusNotFound)
}

func Test_Route_Tokens_Show_Success(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  tokensShowRunner(t, http.StatusOK)
}

//////////////////////////////////////
// CREATE ///////////////////

func tokensCreateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Post("/v1/tokens", AuthorizeOK, binding.Bind(m.TokenAttrs{}), TokensCreate)
  body, _ := json.Marshal(m.TokenAttrs{})
  req, _ := http.NewRequest("POST", "/v1/tokens", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Tokens_Create_Failure_400(t *testing.T) {
  saveToken = func(token *m.Token) (error) {
    token.Errors = []string{"Something went wrong!"}
    return errors.New("*******")
  }
  tokensCreateRunner(t, http.StatusBadRequest)
}

func Test_Route_Tokens_Create_Failure_500(t *testing.T) {
  saveToken = func(token *m.Token) (error) {
    return errors.New("*******")
  }
  tokensCreateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Tokens_Create_Success(t *testing.T) {
  saveToken = func(token *m.Token) (error) {
    return nil
  }
  tokensCreateRunner(t, http.StatusCreated)
}

//////////////////////////////////////
// UPDATE ///////////////////

func tokensUpdateRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Put("/v1/tokens/:token", AuthorizeOK, binding.Bind(m.TokenAttrs{}), TokensUpdate)
  body, _ := json.Marshal(m.TokenAttrs{})
  req, _ := http.NewRequest("PUT", "/v1/tokens/XXXXXX", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Tokens_Update_NotFound(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return nil, errors.New("******")
  }

  tokensUpdateRunner(t, http.StatusNotFound)
}

func Test_Route_Tokens_Update_Failure_400(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  saveToken = func(token *m.Token) (error) {
    token.Errors = []string{"Something went wrong!"}
    return errors.New("*********")
  }

  tokensUpdateRunner(t, http.StatusBadRequest)
}

func Test_Route_Tokens_Update_Failure_500(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  saveToken = func(token *m.Token) (error) {
    return errors.New("*********")
  }

  tokensUpdateRunner(t, http.StatusInternalServerError)
}

func Test_Route_Tokens_Update_Success(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  saveToken = func(token *m.Token) (error) {
    return nil
  }

  tokensUpdateRunner(t, http.StatusOK)
}

//////////////////////////////////////
// DELETE ///////////////////

func tokensDeleteRunner(t *testing.T, code int) {
  server, recorder := testTools(t)
  server.Delete("/v1/tokens/:token", AuthorizeOK, TokensDelete)
  req, _ := http.NewRequest("DELETE", "/v1/tokens/XXXXXX", nil)
  req.Header.Set("Content-Type", "application/json")
  server.ServeHTTP(recorder, req)
  expect(t, recorder.Code, code)
}

func Test_Route_Tokens_Delete_NotFound(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return nil, nil
  }

  tokensDeleteRunner(t, http.StatusNotFound)
}

func Test_Route_Tokens_Delete_Success(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  deleteToken = func(token *m.Token) (error) {
    return nil
  }

  tokensDeleteRunner(t, http.StatusOK)
}

func Test_Route_Tokens_Delete_400(t *testing.T) {
  loadToken = func(id string, u *m.User) (*m.Token, error) {
    expect(t, id, "XXXXXX")
    return &m.Token{}, nil
  }

  deleteToken = func(token *m.Token) (error) {
    return errors.New("*******")
  }

  tokensDeleteRunner(t, http.StatusBadRequest)
}

