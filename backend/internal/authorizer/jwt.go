package authorizer

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"time"

	_ "embed"

	"github.com/go-errors/errors"
	"github.com/go-jose/go-jose/v4"
)

var ErrUnexpectedIssuer = errors.New("unexpected issuer")
var ErrExpiredCredentials = errors.New("expired credentials")
var ErrBadSignature = errors.New("bad signature")
var ErrUnexpectedHTTPStatus = errors.New("unexpected http status")
var ErrEmptyJWKS = errors.New("empty jwks")

//go:embed keys.json
var jwks []byte

type StandardJWTDecoder struct {
	keys jose.JSONWebKeySet
}

const cognitoJWKSURL = "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_Jftnyms2n/.well-known/jwks.json"

func NewStandardJWTDecoder(ctx context.Context) (*StandardJWTDecoder, error) {
	builtInKeys, err := parseJWKS(jwks)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	client := &http.Client{Timeout: 10 * time.Second, Transport: nil, CheckRedirect: nil, Jar: nil}

	keys, err := fetchJWKS(ctx, client, cognitoJWKSURL)
	if err != nil {
		slog.Warn("failed to fetch cognito jwks", "error", err, "action", "falling back to built-in keys")
		keys = builtInKeys
	} else {
		slog.Info("fetched cognito jwks", "url", cognitoJWKSURL, "keys_count", len(keys.Keys))
	}

	return &StandardJWTDecoder{keys: keys}, nil
}

func fetchJWKS(ctx context.Context, client *http.Client, url string) (jose.JSONWebKeySet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return jose.JSONWebKeySet{}, errors.Wrap(err, 0)
	}

	resp, err := client.Do(req)
	if err != nil {
		return jose.JSONWebKeySet{}, errors.Wrap(err, 0)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return jose.JSONWebKeySet{}, errors.Wrap(ErrUnexpectedHTTPStatus, 0)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return jose.JSONWebKeySet{}, errors.Wrap(err, 0)
	}

	return parseJWKS(data)
}

func parseJWKS(data []byte) (jose.JSONWebKeySet, error) {
	var keyList struct {
		Keys []jsontext.Value `json:"keys"`
	}

	if err := json.Unmarshal(data, &keyList); err != nil {
		return jose.JSONWebKeySet{}, errors.Wrap(err, 0)
	}

	var keys jose.JSONWebKeySet

	for _, jsonKey := range keyList.Keys {
		var key jose.JSONWebKey
		if err := key.UnmarshalJSON(jsonKey); err != nil {
			return jose.JSONWebKeySet{}, errors.Wrap(err, 0)
		}

		keys.Keys = append(keys.Keys, key)
	}

	if len(keys.Keys) == 0 {
		return jose.JSONWebKeySet{}, errors.Wrap(ErrEmptyJWKS, 0)
	}

	return keys, nil
}

func (d *StandardJWTDecoder) Decode(jwt string) (Claims, error) {
	signature, err := jose.ParseSigned(jwt, []jose.SignatureAlgorithm{jose.RS256})
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
