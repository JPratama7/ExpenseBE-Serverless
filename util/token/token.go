package token

import (
	"aidanwoods.dev/go-paseto"
	"crud/util/token/option"
	"github.com/JPratama7/util/sync"
	"time"
)

type TokenArgs func(token *paseto.Token) error

type Paseto struct {
	tokenPooler *sync.Pool[*paseto.Token]
	publicKey   paseto.V4AsymmetricPublicKey
	privateKey  paseto.V4AsymmetricSecretKey
	parser      paseto.Parser
	option      option.Option
}

func NewPaseto(publicKey paseto.V4AsymmetricPublicKey, privateKey paseto.V4AsymmetricSecretKey, options ...option.OptionArgs) *Paseto {

	if len(options) == 0 {
		options = []option.OptionArgs{
			option.WithIssuer("default_issuer"),
			option.WithSubject("default_subject"),
			option.WithAudience("default_audience"),
			option.WithExpiration(time.Minute),
		}
	}

	var opts option.Option
	for _, opt := range options {
		opt(&opts)
	}

	return &Paseto{
		tokenPooler: sync.NewPool(func() *paseto.Token {
			token := paseto.NewToken()
			return &token
		}),

		publicKey:  publicKey,
		privateKey: privateKey,
		option:     opts,
		parser:     paseto.NewParser(),
	}
}

func (p *Paseto) Encrypt(options ...TokenArgs) (string, error) {

	token := p.tokenPooler.Get()
	defer p.tokenPooler.Put(token)

	now := time.Now()

	token.SetIssuer(p.option.Issuer)
	token.SetAudience(p.option.Audience)
	token.SetSubject(p.option.Subject)
	token.SetExpiration(now.Add(p.option.Expiration))
	token.SetNotBefore(now)

	for _, opt := range options {
		err := opt(token)
		if err != nil {
			return "", err
		}
	}

	return token.V4Sign(p.privateKey, nil), nil
}

func (p *Paseto) Decrypt(token string) (*paseto.Token, error) {
	return p.parser.ParseV4Public(p.publicKey, token, nil)
}

func WithExpiration(d time.Duration) TokenArgs {
	return func(token *paseto.Token) error {
		token.SetExpiration(time.Now().Add(d))
		return nil
	}
}

func WithNotBefore(d time.Duration) TokenArgs {
	return func(token *paseto.Token) error {
		token.SetNotBefore(time.Now().Add(d))
		return nil
	}
}

func WithIssuer(s string) TokenArgs {
	return func(token *paseto.Token) error {
		token.SetIssuer(s)
		return nil
	}
}

func WithSubject(s string) TokenArgs {
	return func(token *paseto.Token) error {
		token.SetSubject(s)
		return nil
	}
}

func WithAudience(s string) TokenArgs {
	return func(token *paseto.Token) error {
		token.SetAudience(s)
		return nil
	}
}

func WithBody[T any](key string, value T) TokenArgs {
	return func(token *paseto.Token) error {
		return token.Set(key, value)
	}
}

func WithClaims(key string, claims map[string]any) TokenArgs {
	return func(token *paseto.Token) error {
		return token.Set(key, claims)
	}
}

type Token interface {
	Decrypt(token string) (*paseto.Token, error)
	Encrypt(options ...TokenArgs) string
}
