package form3go

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	client = Client{
		PubKeyID:    os.Getenv("FORM3_KEY_ID"),
		PrivKeyPath: os.Getenv("FORM3_PRIV_KEY_PATH"),
		HttpClient:  http.Client{},
	}
)

func TestCreateAccount(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)

	createdAccount, _ := client.CreateAccount(*account)
	assert.NotNil(t, createdAccount)
	assert.Equal(t, *account, createdAccount)
}

func TestFetchAccount(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)

	fetchedAccount, _ := client.FetchAccount("9127e265-9605-4b4b-a0e5-3003ea9cc4dc")
	assert.NotNil(t, fetchedAccount)
	assert.Equal(t, *account, fetchedAccount)

	_, err := client.FetchAccount("")
	assert.Equal(t, "form3go: invalid parameter", err.Error())
}

func TestListAccounts(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)
	accounts := []Account{}
	accounts = append(accounts, *account)

	listedAccounts, _ := client.ListAccounts(0, 1)
	assert.NotNil(t, listedAccounts)
	assert.Equal(t, accounts, listedAccounts)

	listedAccounts, _ = client.ListAccounts(1, 1)
	assert.NotNil(t, listedAccounts)
	assert.Equal(t, []Account{}, listedAccounts)

	listedAccounts, _ = client.ListAccounts(1, 0)
	assert.NotNil(t, listedAccounts)
	assert.Equal(t, []Account{}, listedAccounts)
}

func TestDeleteAccount(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)

	err := client.DeleteAccount("9127e265-9605-4b4b-a0e5-3003ea9cc4dc", "0")
	assert.Nil(t, err)

	err = client.DeleteAccount("9127e265-9605-4b4b-a0e5-3003ea9cc4dc", "2")
	assert.Nil(t, err)

	err = client.DeleteAccount("9127e265-9605-4b4b-a0e5-3003ea9cc4d", "0")
	assert.NotNil(t, err)
}
