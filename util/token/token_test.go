package token

import (
	"aidanwoods.dev/go-paseto"
	"crud/util/token/option"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {

	secretKey := paseto.NewV4AsymmetricSecretKey()

	publicKey := secretKey.Public()

	p := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Issuer = "test_issuer"
		option.Subject = "test_subject"
		option.Audience = "test_audience"
		option.Expiration = time.Hour
	})

	token, err := p.Encrypt(
	//WithBody("key1", "value1"),
	//WithClaims("claims", map[string]any{"claim1": "value1", "claim2": 2}),
	)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if token == "" {
		t.Fatal("Encrypted token is empty")
	}
}

func TestDecrypt(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	publicKey := secretKey.Public()

	p := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Issuer = "test_issuer"
		option.Subject = "test_subject"
		option.Audience = "test_audience"
		option.Expiration = time.Hour
	})

	token, err := p.Encrypt(
		WithBody("key1", "value1"),
		WithClaims("claims", map[string]any{"claim1": "value1", "claim2": 2}),
	)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := p.Decrypt(token)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if is, err := decrypted.GetIssuer(); err != nil || is != "test_issuer" {
		t.Errorf("Expected issuer 'test_issuer', got '%s'", is)
	}
	if sub, err := decrypted.GetSubject(); err != nil || sub != "test_subject" {
		t.Errorf("Expected subject 'test_subject', got '%s'", sub)
	}
	if aud, err := decrypted.GetAudience(); err != nil || aud != "test_audience" {
		t.Errorf("Expected audience 'test_audience', got '%s'", aud)
	}
}

func TestCustomClaims(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	publicKey := secretKey.Public()

	p := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Issuer = "test_issuer"
		option.Subject = "test_subject"
		option.Audience = "test_audience"
		option.Expiration = time.Hour
	})

	token, _ := p.Encrypt(
		WithBody("key1", "value1"),
		WithClaims("claims", map[string]any{"claim1": "value1", "claim2": 2}),
	)

	decrypted, _ := p.Decrypt(token)

	// Check custom claims
	var key1Value string
	err := decrypted.Get("key1", &key1Value)
	if err != nil || key1Value != "value1" {
		t.Errorf("Expected key1 value 'value1', got '%s'", key1Value)
	}

	claims := make(map[string]any, 1)
	err = decrypted.Get("claims", &claims)
	if err != nil {
		t.Errorf("Failed to get claims: %v", err)
	}
	if claims["claim1"] != "value1" || claims["claim2"] != float64(2) {
		t.Errorf("Claims do not match expected values")
	}
}

func TestTokenExpiration(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public()

	shortLivedPaseto := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Expiration = time.Millisecond
	})

	token, err := shortLivedPaseto.Encrypt()
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	time.Sleep(time.Millisecond * 2)

	_, err = shortLivedPaseto.Decrypt(token)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestWithBody(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	publicKey := secretKey.Public()

	p := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Issuer = "test_issuer"
		option.Subject = "test_subject"
		option.Audience = "test_audience"
		option.Expiration = time.Hour
	})

	enc, err := p.Encrypt(WithBody("key", "value"))
	if err != nil {
		t.Fatalf("WithBody failed: %v", err)
	}

	token, err := p.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	var value string
	err = token.Get("key", &value)
	if err != nil || value != "value" {
		t.Errorf("Expected value 'value', got '%s'", value)
	}
}

func TestWithClaims(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	publicKey := secretKey.Public()

	p := NewPaseto(publicKey, secretKey, func(option *option.Option) {
		option.Issuer = "test_issuer"
		option.Subject = "test_subject"
		option.Audience = "test_audience"
		option.Expiration = time.Hour
	})

	claims := map[string]any{"claim1": "value1", "claim2": 2}
	enc, err := p.Encrypt(WithClaims("claims", claims))
	if err != nil {
		t.Fatalf("WithClaims failed: %v", err)
	}

	token, err := p.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	var retrievedClaims map[string]any
	err = token.Get("claims", &retrievedClaims)
	if err != nil {
		t.Fatalf("Failed to get claims: %v", err)
	}
	if retrievedClaims["claim1"] != "value1" || retrievedClaims["claim2"] != float64(2) {
		t.Errorf("Retrieved claims do not match original claims")
	}
}
