package api

import (
  m "github.com/keighl/gypsy-cab/models"
  r "github.com/dancannon/gorethink"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "errors"
)

type TokenData struct {
  User *m.User `json:"current_user,omitempty"`
  *m.Token `json:"token,omitempty"`
}

type TokensData struct {
  User *m.User `json:"current_user,omitempty"`
  Tokens []m.Token `json:"tokens,omitempty"`
}

////////////////////////

var loadTokens = func(u *m.User) ([]m.Token, error) {
  tokens := []m.Token{}
  res, err := r.Table("tokens").
    GetAllByIndex("user_id", u.Id).
    OrderBy(r.Desc("created_at")).
    Run(DB)
  if (err != nil) { return nil, err }
  err = res.All(&tokens)
  return tokens, err
}

func TokensIndex(r render.Render, user *m.User) {
  tokens, err := loadTokens(user)

  if (err != nil) {
    r.JSON(500, MessageEnvelope(err.Error()))
    return
  }

  data := &TokensData{User: user, Tokens: tokens}
  r.JSON(200, data)
}

/////////////////////////

var loadToken = func(id string, u *m.User) (*m.Token, error) {
  token := &m.Token{}
  res, err := r.Table("tokens").Get(id).Run(DB)
  if (err != nil) { return nil, err }
  err = res.One(token)
  if (err != nil) { return nil, err }
  if (token.UserId != u.Id) { return nil, errors.New("Not your token") }
  return token, err
}

func TokensShow(params martini.Params, r render.Render, user *m.User) {
  token, err := loadToken(params["token"], user)

  if (token == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  data := &TokenData{User: user, Token: token}
  r.JSON(200, data)
}

/////////////////////////

var saveToken = func(token *m.Token) (error) {
  return m.Save(token)
}

func TokensCreate(r render.Render, user *m.User, attrs m.TokenAttrs) {
  token := attrs.Token()
  token.UserId = user.Id
  err := saveToken(token)

  if (err != nil) {
    if (token.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), token.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &TokenData{User: user, Token: token}
  r.JSON(201, data)
}

/////////////////////////

func TokensUpdate(params martini.Params, r render.Render, user *m.User, attrs m.TokenAttrs) {
  token, err := loadToken(params["token"], user)

  if (token == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  token.UpdateFromAttrs(attrs)
  err = saveToken(token)

  if (err != nil) {
    if (token.HasErrors()) {
      r.JSON(400, ErrorEnvelope(err.Error(), token.Errors))
    } else {
      r.JSON(500, ServerErrorEnvelope())
    }
    return
  }

  data := &TokenData{User: user, Token: token}
  r.JSON(200, data)
}

/////////////////////////

var deleteToken = func(token *m.Token) (error) {
  return m.Delete(token)
}

func TokensDelete(params martini.Params, r render.Render, user *m.User) {
  token, err := loadToken(params["token"], user)

  if (token == nil || err != nil) {
    r.JSON(404, MessageEnvelope("Record not found"))
    return
  }

  err = deleteToken(token)

  if (err != nil) {
    r.JSON(400, MessageEnvelope("Couldn't delete the item"))
    return
  }

  data := &Data{User: user, Message: &Message{"The token was deleted"}}
  r.JSON(200, data)
}
