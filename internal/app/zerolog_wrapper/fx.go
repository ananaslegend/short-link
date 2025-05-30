package zerolog_wrapper

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func Module() fx.Option {
	return fx.Options(
		// fx.NopLogger,

		fx.Provide(SetupZerolog),
		// fx.WithLogger(WithZerologFx),

		fx.Invoke(func(lc fx.Lifecycle, logger zerolog.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info().Msg("app starting")

					return nil
				},

				OnStop: func(ctx context.Context) error {
					logger.Info().Msg("app stopping")

					return nil
				},
			})
		}),
	)
}

func WithZerologFx(logger zerolog.Logger) fxevent.Logger {
	return NewZerologLogger(logger)
}

type ZerologFx struct {
	Logger zerolog.Logger
}

func NewZerologLogger(logger zerolog.Logger) fxevent.Logger {
	return ZerologFx{Logger: logger.With().Str("component", "fx").Logger()}
}

// LogEvent logs the given event to the provided Logger.
func (l ZerologFx) LogEvent(event fxevent.Event) {
	l.logEvent(event)
}

// logEvent handles different event types and delegates to appropriate handlers.
func (l ZerologFx) logEvent(event fxevent.Event) {
	switch e := event.(type) {
	// Lifecycle events
	case *fxevent.OnStartExecuting, *fxevent.OnStopExecuting, *fxevent.Invoking, *fxevent.Stopping:
		l.logLifecycleEvent(e)

	// Result events
	case *fxevent.OnStartExecuted, *fxevent.OnStopExecuted, *fxevent.Invoked,
		*fxevent.Stopped, *fxevent.RolledBack, *fxevent.Started, *fxevent.LoggerInitialized:
		l.logResultEvent(e)

	// Error events
	case *fxevent.RollingBack:
		l.logRollingBack(e)

	// Dependency events
	case *fxevent.Supplied, *fxevent.Run:
		l.logSimpleDependencyEvent(e)

	// Complex dependency events
	case *fxevent.Provided, *fxevent.Replaced, *fxevent.Decorated:
		l.logComplexDependencyEvent(e)
	}
}

// logLifecycleEvent handles lifecycle events.
func (l ZerologFx) logLifecycleEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logOnStartExecuting(e)
	case *fxevent.OnStopExecuting:
		l.logOnStopExecuting(e)
	case *fxevent.Invoking:
		l.logInvoking(e)
	case *fxevent.Stopping:
		l.logStopping(e)
	}
}

// logResultEvent handles result events.
func (l ZerologFx) logResultEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuted:
		l.logOnStartExecuted(e)
	case *fxevent.OnStopExecuted:
		l.logOnStopExecuted(e)
	case *fxevent.Invoked:
		l.logInvoked(e)
	case *fxevent.Stopped:
		l.logStopped(e)
	case *fxevent.RolledBack:
		l.logRolledBack(e)
	case *fxevent.Started:
		l.logStarted(e)
	case *fxevent.LoggerInitialized:
		l.logLoggerInitialized(e)
	}
}

// logSimpleDependencyEvent handles simple dependency events.
func (l ZerologFx) logSimpleDependencyEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Supplied:
		l.logSupplied(e)
	case *fxevent.Run:
		l.logRun(e)
	}
}

// logComplexDependencyEvent handles complex dependency events.
func (l ZerologFx) logComplexDependencyEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Provided:
		l.logProvided(e)
	case *fxevent.Replaced:
		l.logReplaced(e)
	case *fxevent.Decorated:
		l.logDecorated(e)
	}
}

func (l ZerologFx) logOnStartExecuting(e *fxevent.OnStartExecuting) {
	l.Logger.Info().
		Str("callee", e.FunctionName).
		Str("caller", e.CallerName).
		Msg("OnStart hook executing")
}

func (l ZerologFx) logOnStartExecuted(e *fxevent.OnStartExecuted) {
	logger := l.Logger.Info()
	if e.Err != nil {
		logger = l.Logger.Error().Err(e.Err)
		logger.Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook failed")
	} else {
		logger.Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Str("runtime", e.Runtime.String()).
			Msg("OnStart hook executed")
	}
}

func (l ZerologFx) logOnStopExecuting(e *fxevent.OnStopExecuting) {
	l.Logger.Info().
		Str("callee", e.FunctionName).
		Str("caller", e.CallerName).
		Msg("OnStop hook executing")
}

func (l ZerologFx) logOnStopExecuted(e *fxevent.OnStopExecuted) {
	logger := l.Logger.Info()
	if e.Err != nil {
		logger = l.Logger.Error().Err(e.Err)
		logger.Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook failed")
	} else {
		logger.Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Str("runtime", e.Runtime.String()).
			Msg("OnStop hook executed")
	}
}

func (l ZerologFx) logSupplied(e *fxevent.Supplied) {
	if e.Err != nil {
		l.Logger.Error().Err(e.Err).Msg("supplied error")
	} else {
		l.Logger.Info().
			Str("type", e.TypeName).
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Msg("supplied")
	}
}

func (l ZerologFx) logProvided(e *fxevent.Provided) {
	for _, rtype := range e.OutputTypeNames {
		l.Logger.Info().
			Str("constructor", e.ConstructorName).
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Str("type", rtype).
			Bool("private", e.Private).
			Msg("provided")
	}

	if e.Err != nil {
		l.Logger.Error().
			Str("module", e.ModuleName).
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Err(e.Err).
			Msg("error encountered while applying options")
	}
}

func (l ZerologFx) logReplaced(e *fxevent.Replaced) {
	for _, rtype := range e.OutputTypeNames {
		l.Logger.Info().
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Str("type", rtype).
			Msg("replaced")
	}

	if e.Err != nil {
		l.Logger.Error().
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Err(e.Err).
			Msg("error encountered while replacing")
	}
}

func (l ZerologFx) logDecorated(e *fxevent.Decorated) {
	for _, rtype := range e.OutputTypeNames {
		l.Logger.Info().
			Str("decorator", e.DecoratorName).
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Str("type", rtype).
			Msg("decorated")
	}

	if e.Err != nil {
		l.Logger.Error().
			Strs("stacktrace", e.StackTrace).
			Strs("moduletrace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Err(e.Err).
			Msg("error encountered while applying options")
	}
}

func (l ZerologFx) logRun(e *fxevent.Run) {
	if e.Err != nil {
		l.Logger.Error().Err(e.Err).
			Str("name", e.Name).
			Str("kind", e.Kind).
			Str("module", e.ModuleName).
			Msg("run error")
	} else {
		l.Logger.Info().
			Str("name", e.Name).
			Str("kind", e.Kind).
			Str("runtime", e.Runtime.String()).
			Str("module", e.ModuleName).
			Msg("run")
	}
}

func (l ZerologFx) logInvoking(e *fxevent.Invoking) {
	l.Logger.Info().
		Str("function", e.FunctionName).
		Str("module", e.ModuleName).
		Msg("invoking")
}

func (l ZerologFx) logInvoked(e *fxevent.Invoked) {
	if e.Err != nil {
		l.Logger.Error().
			Err(e.Err).
			Str("stack", e.Trace).
			Str("function", e.FunctionName).
			Str("module", e.ModuleName).
			Msg("invoke failed")
	}
}

func (l ZerologFx) logStopping(e *fxevent.Stopping) {
	l.Logger.Info().
		Str("signal", strings.ToUpper(e.Signal.String())).
		Msg("received signal")
}

func (l ZerologFx) logStopped(e *fxevent.Stopped) {
	if e.Err != nil {
		l.Logger.Error().
			Err(e.Err).
			Msg("stop failed")
	}
}

func (l ZerologFx) logRollingBack(e *fxevent.RollingBack) {
	l.Logger.Error().
		Err(e.StartErr).
		Msg("start failed, rolling back")
}

func (l ZerologFx) logRolledBack(e *fxevent.RolledBack) {
	if e.Err != nil {
		l.Logger.Error().
			Err(e.Err).
			Msg("rollback failed")
	}
}

func (l ZerologFx) logStarted(e *fxevent.Started) {
	if e.Err != nil {
		l.Logger.Error().
			Err(e.Err).
			Msg("start failed")
	} else {
		l.Logger.Info().Msg("started")
	}
}

func (l ZerologFx) logLoggerInitialized(e *fxevent.LoggerInitialized) {
	if e.Err != nil {
		l.Logger.Error().
			Err(e.Err).
			Msg("custom Logger initialization failed")
	} else {
		l.Logger.Info().
			Str("function", e.ConstructorName).
			Msg("initialized custom fxevent.Logger")
	}
}
