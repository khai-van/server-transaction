package api

import (
	"fmt"
	"time"

	"app/pkg/database"

	"github.com/kataras/iris/v12"
)

type Server struct {
	*iris.Application
	config Configuration

	db *database.DB
}

func NewServer(c Configuration) *Server {
	app := iris.Default().SetName(c.ServerName)
	app.Configure(iris.WithLowercaseRouting)

	srv := &Server{
		Application: app,
		config:      c,
		db:          &database.DB{},
	}
	//init database mongo
	srv.db.InitMongoDB(srv.config.MongoConfig.DSN, srv.config.MongoConfig.DB)

	srv.buildRouter()

	return srv
}

func (srv *Server) Start() error {
	// config timeout
	srv.ConfigureHost(func(su *iris.Supervisor) {
		su.Server.ReadTimeout = time.Minute
		su.Server.WriteTimeout = time.Minute
		su.Server.IdleTimeout = time.Minute
		su.Server.ReadHeaderTimeout = time.Minute
	})

	addr := fmt.Sprintf("%s:%d", srv.config.Host, srv.config.Port)
	return srv.Listen(addr)
}
