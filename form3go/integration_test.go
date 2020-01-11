// +build integration

package form3go

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)

	// Create account
	client := Client{
		PubKeyID:    os.Getenv("FORM3_KEY_ID"),
		PrivKeyPath: os.Getenv("FORM3_PRIV_KEY_PATH"),
		HttpClient:  http.Client{},
	}
	acct, err := client.CreateAccount(*account)
	assert.Nil(t, err)
	assert.Equal(t, *account, acct)

	// Fetch account
	acct, err = client.FetchAccount(account.AccountData.ID)
	assert.Nil(t, err)
	assert.Equal(t, *account, acct)

	// List accounts
	accts, err := client.ListAccounts(0, 1)
	accounts := []Account{}
	accounts = append(accounts, *account)
	assert.Nil(t, err)
	assert.Equal(t, accounts, accts)

	// Delete account
	err = client.DeleteAccount(account.AccountData.ID, "0")
	assert.Nil(t, err)
}
