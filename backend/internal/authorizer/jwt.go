package authorizer

import (
	"encoding/json"
	"time"

	_ "embed"

	"github.com/go-errors/errors"
	"gopkg.in/square/go-jose.v2"
)

var ErrUnexpectedIssuer = errors.New("unexpected issuer")
var ErrExpiredCredentials = errors.New("expired credentials")
var ErrBadSignature = errors.New("bad signature")

//go:embed keys.json
var jwks []byte

type StandardJWTDecoder struct {
	keys jose.JSONWebKeySet
}

func NewStandardJWTDecoder() (*StandardJWTDecoder, error) {
	var keyList struct {
		Keys []json.RawMessage `json:"keys"`
	}

	err := json.Unmarshal(jwks, &keyList)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	var keys jose.JSONWebKeySet

	for _, jsonKey := range keyList.Keys {
		k := jose.JSONWebKey{}
		if err := k.UnmarshalJSON(jsonKey); err != nil {
			return nil, errors.Wrap(err, 0)
		}

		keys.Keys = append(keys.Keys, k)
	}

	return &StandardJWTDecoder{
		keys: keys,
	}, nil
}

func (d *StandardJWTDecoder) Decode(jwt string) (Claims, error) {
	signature, err := jose.ParseSigned(jwt)
	if err != nil {
		return Claims{}, errors.Wrap(err, 0)
	}

	kid := signature.Signatures[0].Header.KeyID
	var key interface{}
	if result := d.keys.Key(kid); len(result) == 1 {
		key = result[0].Key
	} else {
		return Claims{}, errors.Wrap(ErrUnexpectedIssuer, 0)
	}

	payload, err := signature.Verify(key)
	if err != nil {
		return Claims{}, errors.Wrap(ErrBadSignature, 0)
	}

	var claims Claims

	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, errors.Wrap(err, 0)
	}

	if time.Unix(claims.Expiration, 0).Before(time.Now()) {
		return Claims{}, errors.Wrap(ErrExpiredCredentials, 0)
	}

	return claims, nil
}
