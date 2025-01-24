package apis

import (
	"github.com/nanoteck137/beldum"
	"github.com/nanoteck137/beldum/core"
	"github.com/nanoteck137/pyrin"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: beldum.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	return s, nil
}
