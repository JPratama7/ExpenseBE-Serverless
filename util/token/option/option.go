package option

import "time"

type Option struct {
	Issuer     string
	Subject    string
	Audience   string
	Expiration time.Duration
}

type OptionArgs func(option *Option)

func WithIssuer(s string) OptionArgs {
	return func(option *Option) {
		option.Issuer = s
	}
}

func WithSubject(s string) OptionArgs {
	return func(option *Option) {
		option.Subject = s
	}
}

func WithAudience(s string) OptionArgs {
	return func(option *Option) {
		option.Audience = s
	}
}

func WithExpiration(d time.Duration) OptionArgs {
	return func(option *Option) {
		option.Expiration = d
	}
}
