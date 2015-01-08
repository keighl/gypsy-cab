package main

import (
  "gypsy/api"
  m "gypsy/models"
  u "gypsy/utils"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
)

/////////////////////////////

func main() {
  config := u.ConfigForFile("config/app.json")
  DB := u.RethinkSession(config)
  api.DB = DB
  api.Config = u.ConfigForFile("config/app.json")
  m.DB = DB
  m.Config = u.ConfigForFile("config/app.json")
  server := u.MartiniServer(config.ServerLoggingEnabled)
  SetupServerRoutes(server)
  server.Run() // Blocks....
  DB.Close()
}

func SetupServerRoutes(server *martini.ClassicMartini) {

  server.Get("/v1/", api.Authorize, api.Me)

  // Login
  server.Post("/v1/login", binding.Bind(m.UserAttrs{}), api.Login)

  // Signup
  server.Post("/v1/users", binding.Bind(m.UserAttrs{}), api.UserCreate)

  // Me
  server.Get("/v1/users/me", api.Authorize, api.Me)
  server.Put("/v1/users/me", api.Authorize, binding.Bind(m.UserAttrs{}), api.MeUpdate)

  // Jobs
  server.Get("/v1/jobs", api.Authorize, api.JobsIndex)
  server.Post("/v1/jobs", api.Authorize, binding.Bind(m.JobAttrs{}), api.JobsCreate)
  server.Get("/v1/jobs/:job_id", api.Authorize, api.JobsShow)
  server.Put("/v1/jobs/:job_id", api.Authorize, binding.Bind(m.JobAttrs{}), api.JobsUpdate)
  server.Delete("/v1/jobs/:job_id", api.Authorize, api.JobsDelete)
  server.Post("/v1/jobs/:job_id/:token", binding.Bind(api.Upload{}), api.JobProcess)
  server.Get("/v1/jobs/:job_id/items", api.Authorize, api.ItemsIndex)

  // Tokens
  server.Get("/v1/tokens", api.Authorize, api.TokensIndex)
  server.Post("/v1/tokens", api.Authorize, binding.Bind(m.TokenAttrs{}), api.TokensCreate)
  server.Get("/v1/tokens/:token", api.Authorize, api.TokensShow)
  server.Put("/v1/tokens/:token", api.Authorize, binding.Bind(m.TokenAttrs{}), api.TokensUpdate)
  server.Delete("/v1/tokens/:token", api.Authorize, api.TokensDelete)

  // Password Reset
  server.Post("/v1/password-reset", binding.Bind(m.PasswordResetAttrs{}), api.PasswordResetCreate)
  server.Post("/v1/password-reset/:token", binding.Bind(m.UserAttrs{}), api.PasswordResetUpdate)
}


