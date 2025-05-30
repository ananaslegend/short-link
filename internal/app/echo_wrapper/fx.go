package echowraper

import (
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/echo_wrapper/mw"
)

func Module() fx.Option {
	return fx.Module(
		"shortlink.internal.app.echo",
		fx.Provide(NewEchoRouter),
		fx.Provide(NewEchoAPIGroup),

		fx.Invoke(SetupValidator),
		fx.Invoke(mw.SetupMiddleware),
		fx.Invoke(RunEchoServer),
	)
}
