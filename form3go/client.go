package form3go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var (
	host    = os.Getenv("FORM3_HOST")
	baseURL = "http://" + host

	acctURL    = "/v1/organisation/accounts"
	acctReqURL = baseURL + acctURL

	// Errors used by the library

	// ErrInvalidAccount is returned by CreateAccount when account
	// information is invalid.
	ErrInvalidAccount = errors.New("form3go: invalid request body")

	// ErrEmptyHost is returned by CreateAccount when FORM3_HOST env
	// variable is not provided.
	ErrEmptyHost = errors.New("form3go: FORM3_HOST env variable is required")

	// ErrParameterEmpty is returned by FetchAccount and DeleteAccount
	// ListAccounts when provided parameters are empty.
	ErrParameterEmpty = errors.New("form3go: invalid parameter")

	// ErrCreateAccount is returned by CreateAccount when creating
	// account is failed.
	ErrCreateAccount = errors.New("form3go: create account failure")

	// ErrDeleteAccount is returned by DeleteAccount when deleting
	// account is failed.
	ErrDeleteAccount = errors.New("form3go: delete account failure")
)

type Client struct {
	PubKeyID    string
	PrivKeyPath string

	HttpClient http.Client
}

// CreateAccount creates account.
func (c *Client) CreateAccount(acct Account) (Account, error) {
	// validate given account info
	if err := acct.Validate(); err != nil {
		return Account{}, ErrInvalidAccount
	}

	// create request
	acctByte, err := json.Marshal(acct)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected JSON marshal failure: %v", err)
	}
	req, err := http.NewRequest("POST", acctReqURL, bytes.NewBuffer(acctByte))
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}

	// generate header informations
	reqInfo := request{
		data:     string(acctByte),
		endpoint: acctReqURL,
		method:   "POST",
		keyPath:  c.PrivKeyPath,
		keyID:    c.PubKeyID,
	}
	date := reqInfo.genDateHeader()
	digest := reqInfo.genDigestHeader()
	sig, err := reqInfo.genSignature(date, digest)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected generating Signature failure: %v", err)
	}
	authHeader, err := reqInfo.genAuthHeader(sig)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected generating Authorization header failure: %v", err)
	}
	req.Header.Set("Host", os.Getenv("FORM3_HOST"))
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Digest", digest)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Content-Length", strconv.Itoa(len(string(acctByte))))

	// do request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode != 201 {
		return Account{}, ErrCreateAccount
	}
	account := Account{}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected response decode failure: %v", err)
	}

	return account, nil
}

// FetchAccount fetches account with ID
func (c *Client) FetchAccount(id string) (Account, error) {
	// check id
	if id == "" {
		return Account{}, ErrParameterEmpty
	}

	// create requet
	req, err := http.NewRequest("GET", acctReqURL+"/"+id, nil)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}

	// generate headers
	reqInfo := request{
		endpoint: acctReqURL,
		method:   "GET",
		keyPath:  c.PrivKeyPath,
		keyID:    c.PubKeyID,
	}
	date := reqInfo.genDateHeader()
	sig, err := reqInfo.genSignature(date, "")
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected generating Signature failure: %v", err)
	}
	authHeader, err := reqInfo.genAuthHeader(sig)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected generating Authorization header failure: %v", err)
	}
	req.Header.Set("Host", os.Getenv("FORM3_HOST"))
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)

	// do request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}
	defer resp.Body.Close()

	// check response
	account := Account{}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return Account{}, fmt.Errorf("form3go: unexpected response decode failure: %v", err)
	}

	return account, nil
}

// ListAccounts returns array of accounts
func (c *Client) ListAccounts(pageNumber, pageSize int) ([]Account, error) {
	// create request
	num := strconv.Itoa(pageNumber)
	size := strconv.Itoa(pageSize)
	req, err := http.NewRequest("GET", acctReqURL+"?page[number]="+num+"&page[size]="+size, nil)
	if err != nil {
		return []Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}

	// generate headers
	reqInfo := request{
		endpoint: acctReqURL,
		method:   "GET",
		keyPath:  c.PrivKeyPath,
		keyID:    c.PubKeyID,
	}
	date := reqInfo.genDateHeader()
	sig, err := reqInfo.genSignature(date, "")
	if err != nil {
		return []Account{}, fmt.Errorf("form3go: unexpected generating Signature failure: %v", err)
	}
	authHeader, err := reqInfo.genAuthHeader(sig)
	if err != nil {
		return []Account{}, fmt.Errorf("form3go: unexpected generating Authorization header failure: %v", err)
	}
	req.Header.Set("Host", os.Getenv("FORM3_HOST"))
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)

	// do request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return []Account{}, fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}
	defer resp.Body.Close()

	// check response
	accts := struct {
		Accounts []Data `json:"data"`
	}{
		Accounts: []Data{},
	}
	err = json.NewDecoder(resp.Body).Decode(&accts)
	if err != nil {
		return []Account{}, fmt.Errorf("form3go: unexpected response decode failure: %v", err)
	}

	// adjust response
	accounts := []Account{}
	for _, v := range accts.Accounts {
		account := Account{
			AccountData: v,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// DeleteAccount removes account with ID
func (c *Client) DeleteAccount(id, version string) error {
	// check parameters
	if id == "" || version == "" {
		return ErrParameterEmpty
	}

	// create request
	req, err := http.NewRequest("DELETE", acctReqURL+"/"+id+"?version="+version, nil)
	if err != nil {
		return fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}

	// create request
	reqInfo := request{
		endpoint: acctReqURL,
		method:   "DELETE",
		keyPath:  c.PrivKeyPath,
		keyID:    c.PubKeyID,
	}
	date := reqInfo.genDateHeader()
	sig, err := reqInfo.genSignature(date, "")
	if err != nil {
		return fmt.Errorf("form3go: unexpected generating Signature failure: %v", err)
	}
	authHeader, err := reqInfo.genAuthHeader(sig)
	if err != nil {
		return fmt.Errorf("form3go: unexpected generating Authorization header failure: %v", err)
	}
	req.Header.Set("Host", os.Getenv("FORM3_HOST"))
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)

	// do request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("form3go: unexpected HTTP request failure: %v", err)
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode != 204 {
		return ErrDeleteAccount
	}

	return nil
}
