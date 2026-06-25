package sqlite

import "time"

// Option configures an opened SQLite database.
type Option func(*options)

type options struct {
	busyTimeout time.Duration
}

// WithBusyTimeout sets the SQLite busy timeout for lock contention.
func WithBusyTimeout(d time.Duration) Option {
	return func(o *options) {
		o.busyTimeout = d
	}
}

func defaultOptions() options {
	return options{
		busyTimeout: 5 * time.Second,
	}
}
