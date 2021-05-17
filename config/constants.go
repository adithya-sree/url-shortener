package config

const (
	ConfigPath = "/Users/exerdra/git/url-shortener/sample-config.json"

	ErrString = "error"

	UrlDoesNotExist       = "requested url does not exist"
	ErrorRedirecting      = "error redirecting url"
	ErrorCreatingRedirect = "error creating redirect"

	Session = "session-id"
	Url     = "url"

	HeaderUrl     = "x-shortener-url"
	HeaderSession = "x-session-id"
)
