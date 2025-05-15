package session

import "time"

type optionFn func(*sessionManagerConfig)

func WithSessionName(sessionName string) optionFn {
	return func(conf *sessionManagerConfig) {
		conf.sessionName = sessionName
	}
}
func WithCsrfFormField(csrfFormField string) optionFn {
	return func(conf *sessionManagerConfig) {
		conf.csrfFormField = csrfFormField
	}
}
func WithUnAuthLifeTime(unAuthLifeTime time.Duration) optionFn {
	return func(conf *sessionManagerConfig) {
		conf.unAuthLifeTime = unAuthLifeTime
	}
}
func WithAuthLifeTime(authLifeTime time.Duration) optionFn {
	return func(conf *sessionManagerConfig) {
		conf.authLifeTime = authLifeTime
	}
}
func WithIdleTimeout(idleTimeout time.Duration) optionFn {
	return func(conf *sessionManagerConfig) {
		conf.idleTimeout = idleTimeout
	}
}
