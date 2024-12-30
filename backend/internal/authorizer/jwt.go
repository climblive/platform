package authorizer

import (
	"encoding/json"
	"fmt"
	"time"

	_ "embed"

	"github.com/go-errors/errors"
	"gopkg.in/square/go-jose.v2"
)

var ErrUnexpectedIssuer = errors.New("unexpected issuer")
var ErrExpiredCredentials = errors.New("expired credentials")

type claims struct {
	Username   string `json:"username"`
	Expiration int64  `json:"exp"`
	Issuer     string `json:"iss"`
}

var keys jose.JSONWebKeySet

//go:embed keys.json
var jwks []byte

func init() {
	var keyList struct {
		Keys []json.RawMessage `json:"keys"`
	}

	err := json.Unmarshal(jwks, &keyList)
	if err != nil {
		panic(fmt.Errorf("unmarshal jwks: %w", err))
	}

	for _, jsonKey := range keyList.Keys {
		k := jose.JSONWebKey{}
		if err := k.UnmarshalJSON(jsonKey); err != nil {
			panic(fmt.Errorf("unmarshal jwk: %w", err))
		}

		keys.Keys = append(keys.Keys, k)
	}
}

func (a *Authorizer) verifyJWT(jwt string) ([]byte, error) {
	signature, err := jose.ParseSigned(jwt)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	kid := signature.Signatures[0].Header.KeyID
	var key interface{}
	if result := keys.Key(kid); len(result) == 1 {
		key = result[0].Key
	} else {
		return nil, ErrUnexpectedIssuer
	}

	payload, err := signature.Verify(key)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return payload, nil
}

func (a *Authorizer) decodeJWT(payload []byte) (string, error) {
	var claims claims

	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}

	if time.Unix(claims.Expiration, 0).Before(time.Now()) {
		return "", ErrExpiredCredentials
	}

	return claims.Username, nil
}
