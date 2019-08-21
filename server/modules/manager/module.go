package manager

type startupFunc func() error
type shutdownFunc func()

var startups []startupFunc
var shutdowns []shutdownFunc

func OnStartup(f startupFunc) {
	startups = append(startups, f)
}

func OnShutdown(f shutdownFunc) {
	shutdowns = append(shutdowns, f)
}

func Startup() error {
	for _, f := range startups {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func Shutdown() {
	for _, f := range shutdowns {
		f()
	}
}
