package client

import (
	"crypto/hmac"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

const (
	oauthSignatureMethod = "HMAC-SHA1"
	oauthVersion         = "1.0"
	oauthNonceLength     = 32
	authKeysLength       = 7
)

const (
	oauthConsumerKey        = "oauth_consumer_key"
	oauthTokenKey           = "oauth_token"
	oauthSignatureMethodKey = "oauth_signature_method"
	oauthTimestampKey       = "oauth_timestamp"
	oauthNonceKey           = "oauth_nonce"
	oauthVersionKey         = "oauth_version"
	oauthSignatureKey       = "oauth_signature"
	oAuthKey                = "OAuth "
	statusKey               = "status"
)

type authHandler struct {
	oAuthSignatureMethod string
	oAuthVersion         string
	authKeys             [authKeysLength]string
}

func NewAuthHandler() SignatureCreator {
	return authHandler{
		oAuthSignatureMethod: oauthSignatureMethod,
		oAuthVersion:         oauthVersion,
		authKeys: [authKeysLength]string{oauthConsumerKey, oauthTokenKey, oauthSignatureMethodKey, oauthTimestampKey,
			oauthNonceKey, oauthVersionKey, oauthSignatureKey},
	}
}

// CreateAuthorizationString - To build the header string, imagine writing to a string named DST.
// Append the string “OAuth ” (including the space at the end) to DST.
// For each key/value pair of the 7 parameters listed above:
// Percent encode the key and append it to DST.
// Append the equals character ‘=’ to DST.
// Append a double quote ‘”’ to DST.
// Percent encode the value and append it to DST.
// Append a double quote ‘”’ to DST.
// If there are key/value pairs remaining, append a comma ‘,’ and a space ‘ ‘ to DST.
func (a authHandler) CreateAuthorizationString(s SignatureRequest) string {
	oauthTimestamp := createTimestamp()
	base64OauthNonce := encodeToB64String(randStringBytes(oauthNonceLength))

	signature := a.createSignature(s, base64OauthNonce, oauthTimestamp)

	authorizationValues := [authKeysLength]string{s.key, s.token, a.oAuthSignatureMethod,
		oauthTimestamp, base64OauthNonce, a.oAuthVersion, signature}

	authStr := oAuthKey
	for iter, key := range a.authKeys {
		authStr += url.QueryEscape(key) + "=\"" + url.QueryEscape(authorizationValues[iter]) + "\""
		if iter != authKeysLength-1 {
			authStr += ","
		}
	}

	return authStr
}

func createTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func encodeToB64String(s string) string {
	return b64.StdEncoding.EncodeToString([]byte(s))
}

// The base URL is the URL to which the request is directed, minus any query string or hash parameters
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/creating-a-signature
func (a authHandler) createSignature(s SignatureRequest, base64OauthNonce, oauthTimestamp string) string {

	s.params.Add(oauthConsumerKey, s.key)
	s.params.Add(oauthNonceKey, base64OauthNonce)
	s.params.Add(oauthSignatureMethodKey, a.oAuthSignatureMethod)
	s.params.Add(oauthTimestampKey, oauthTimestamp)
	s.params.Add(oauthTokenKey, s.token)
	s.params.Add(oauthVersionKey, a.oAuthVersion)
	if s.body != "" {
		s.params.Add(statusKey, s.body) // request body (POST method)
	}

	escapedStr := fmt.Sprintf("%s&%s&%s", s.method, url.QueryEscape(s.url), url.QueryEscape(s.params.Encode()))
	signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(s.secret), url.QueryEscape(s.accessSecret))
	signature := makeSignature(escapedStr, signingKey)

	return signature
}

// makeSignature - create final value by passing the signature base string and signing key
// to the HMAC-SHA1 hashing algorithm
func makeSignature(input, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(input))

	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}

// RandStringBytes - make unique random string for each request (oauth_nonce parameter)
func randStringBytes(stringLength int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, stringLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
