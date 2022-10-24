package app

// BootLoader for app start boot.
type BootLoader interface {
	// Boot do something before application run
	Boot(ka *KiteApp) error
}

// BootChecker for app start boot.
type BootChecker interface {
	// BeforeBoot check something before boot run
	BeforeBoot() bool
}

// BootFunc for application
type BootFunc func(ka *KiteApp) error

// Boot do something
func (fn BootFunc) Boot(ka *KiteApp) error {
	return fn(ka)
}

// StdLoader struct
type StdLoader struct {
	BeforeFn func() bool
	BootFn   BootFunc
}

func (l *StdLoader) BeforeBoot() bool {
	if l.BeforeFn != nil {
		return l.BeforeFn()
	}
	return true
}

func (l *StdLoader) Boot(ka *KiteApp) error {
	return l.BootFn(ka)
}
