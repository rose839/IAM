package shutdown

type SMShutdownStartFunc func() error

func (f SMShutdownStartFunc) GetName() string {
	return "test-sm"
}

func (f SMShutdownStartFunc) ShutdownStart() error {
	return f()
}

func (f SMShutdownStartFunc) ShutdownFinish() error {
	return nil
}

func (f SMShutdownStartFunc) Start(gs GSInterface) error {
	return nil
}

type SMFinishFunc func() error

func (f SMFinishFunc) GetName() string {
	return "test-sm"
}

func (f SMFinishFunc) ShutdownStart() error {
	return nil
}

func (f SMFinishFunc) ShutdownFinish() error {
	return f()
}

func (f SMFinishFunc) Start(gs GSInterface) error {
	return nil
}

type SMStartFunc func() error

func (f SMStartFunc) GetName() string {
	return "test-sm"
}

func (f SMStartFunc) ShutdownStart() error {
	return nil
}

func (f SMStartFunc) ShutdownFinish() error {
	return nil
}

func (f SMStartFunc) Start(gs GSInterface) error {
	return f()
}
