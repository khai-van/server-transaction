package api

import (

	"app/api/users"
	"app/transaction"

	"github.com/kataras/iris/v12/middleware/modrevision"
)

func (srv *Server) buildRouter() {
	// Add a health route.
	srv.Get("/health", modrevision.New(modrevision.Options{
		ServerName: srv.config.ServerName,
		Env:        srv.config.Env,
		Developer:  "khaivan",
	}))

	// add main router
	api := srv.Party("/")
	api.RegisterDependency(
		srv.db,
		transaction.NewRepository,
	)

	api.PartyConfigure("/", new(users.UserService))
}
