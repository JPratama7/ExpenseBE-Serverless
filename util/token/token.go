package token

import (
	"aidanwoods.dev/go-paseto"
	"github.com/JPratama7/util/sync"
	"time"
)

type Option struct {
	Issuer     string
	Subject    string
	Audience   string
	Expiration time.Duration
}

type TokenArgs func(token *paseto.Token) error

type OptionArgs func(option *Option)

type Paseto struct {
	tokenPooler *sync.Pool[*paseto.Token]
	publicKey   paseto.V4AsymmetricPublicKey
	privateKey  paseto.V4AsymmetricSecretKey
	parser      paseto.Parser
	option      Option
}

func NewPaseto(publicKey paseto.V4AsymmetricPublicKey, privateKey paseto.V4AsymmetricSecretKey, options ...OptionArgs) *Paseto {

	var option Option
	for _, opt := range options {
		opt(&option)
	}

	return &Paseto{
		tokenPooler: sync.NewPool(func() *paseto.Token {
			return new(paseto.Token)
		}),

		publicKey:  publicKey,
		privateKey: privateKey,
		option:     option,
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

	for _, option := range options {
		err := option(token)
		if err != nil {
			return "", err
		}
	}

	return token.V4Sign(p.privateKey, nil), nil
}

func (p *Paseto) Decrypt(token string) (*paseto.Token, error) {
	return p.parser.ParseV4Public(p.publicKey, token, nil)
}

func WithBody[T any](key string, value T) TokenArgs {
	return func(token *paseto.Token) {
		token.Set(key, value)
	}
}

func WithClaims(key string, claims map[string]any) TokenArgs {
	return func(token *paseto.Token) {
		token.Set(key, claims)
	}
}

type Token interface {
	Decrypt(token string) (*paseto.Token, error)
	Encrypt(options ...TokenArgs) string
}
