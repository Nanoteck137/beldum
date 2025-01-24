package apis

import (
	"github.com/nanoteck137/beldum/core"
	"github.com/nanoteck137/pyrin"
)

func InstallHandlers(app core.App, g pyrin.Group) {
	InstallTaskHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallSystemHandlers(app, g)
	InstallUserHandlers(app, g)
}
