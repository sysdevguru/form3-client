package form3go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// request represents request info to endpoints
type request struct {
	data     string
	endpoint string
	method   string
	keyPath  string
	keyID    string
}

// generate Date Header
func (r *request) genDateHeader() string {
	return time.Now().Format(time.RFC1123)
}

// generate Digest Header
func (r *request) genDigestHeader() string {
	hash := sha256.Sum256([]byte(r.data))
	return "SHA-256=" + string(base64.StdEncoding.EncodeToString(hash[:]))
}

// generate Signature
func (r *request) genSignature(date, digest string) (string, error) {
	if date == "" {
		return "", errors.New("empty date")
	}
	if os.Getenv("FORM3_HOST") == "" {
		return "", errors.New("empty FORM3_HOST env variable")
	}
	if r.keyID == "" {
		return "", errors.New("empty FORM3_KEY_ID env variable")
	}
	if r.keyPath == "" {
		return "", errors.New("empty FORM3_PRIV_KEY_PATH env variable")
	}

	signatureStr := "(request-target): " + r.method + " " + r.endpoint + "\n"
	signatureStr += "host: " + os.Getenv("FORM3_HOST") + "\n"
	signatureStr += "date: " + date + "\n"
	if digest != "" {
		signatureStr += "accept: application/vnd.api+json\n"
		signatureStr += "content-type: application/vnd.api+json\n"
		signatureStr += "content-length: " + strconv.Itoa(len(r.data)) + "\n"
		signatureStr += "digest: " + digest + "\n"
	}

	signer, err := loadPrivateKey(r.keyPath)
	if err != nil {
		return "", err
	}
	signed, err := signer.Sign([]byte(signatureStr))
	if err != nil {
		return "", err
	}
	sig := base64.StdEncoding.EncodeToString(signed)
	return sig, nil
}

// generate Authorization Header
func (r *request) genAuthHeader(sig string) (string, error) {
	if sig == "" {
		return "", errors.New("genAuthHeader: invalid signature")
	}
	return `Signature keyId="` + r.keyID + `",algorithm="rsa-sha256",header="(request-target) host date accept content-type content-length digest",signature="` + sig + `"`, nil
}

func loadPrivateKey(path string) (Signer, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(dat)
}

func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("rsa: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("rsa: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

type Signer interface {
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var rsaKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		rsaKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("rsa: unsupported key type %T", k)
	}
	return rsaKey, nil
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
