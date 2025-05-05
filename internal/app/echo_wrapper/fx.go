package echowraper

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"shortlink.internal.app.echo",
		fx.Provide(NewEchoRouter),
		fx.Provide(NewEchoAPIGroup),

		fx.Invoke(SetupValidator),
		fx.Invoke(SetupMiddleware),
		fx.Invoke(RunEchoServer),
	)
}
