package form3go

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testAccountInfo = `{
		"data": {
		  "type": "accounts",
		  "id": "9127e265-9605-4b4b-a0e5-3003ea9cc4dc",
		  "organisation_id": "db0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		  "attributes": {
			"country": "GB",
			"base_currency": "GBP",
			"account_number": "G1426819",
			"bank_id": "D00300",
			"bank_id_code": "DBDSC",
			"bic": "NWBKGB22",
			"iban": "GB11NWBK40030041426819",
			"title": "23fdc&",
			"first_name": "Samantha",
			"bank_account_name": "7amantha Holder",
			"alternative_bank_account_names": [
				"Sam Holder"
			],
			"account_classification": "Persoasdsdfsd234234",
			"joint_account": false,
			"account_matching_opt_out": false,
			"secondary_identification": "x1B2C3D4"
		  }
		}
	  }`
)

func TestAccountValidate(t *testing.T) {
	account := &Account{}
	_ = json.Unmarshal([]byte(testAccountInfo), account)

	// Empty ID is provided
	account.AccountData.ID = ""
	assert.Equal(t, "Key: 'Account.AccountData.ID' Error:Field validation for 'ID' failed on the 'id' tag", account.Validate().Error())

	// Invalid UUID is provided for ID
	account.AccountData.ID = "127e265-9605-4b4b-a0e5-3003ea9cc4dc"
	assert.Equal(t, "Key: 'Account.AccountData.ID' Error:Field validation for 'ID' failed on the 'id' tag", account.Validate().Error())

	// Invalid Country Code is provided
	account.AccountData.ID = "9127e265-9605-4b4b-a0e5-3003ea9cc4dc"
	account.AccountData.Attributes.Country = "342"
	assert.Equal(t, "Key: 'Account.AccountData.Attributes.Country' Error:Field validation for 'Country' failed on the 'country' tag", account.Validate().Error())

	// Valid Country Code and Account Number
	account.AccountData.Attributes.Country = "GB"
	account.AccountData.Attributes.AccountNumber = "51426819"
	assert.Nil(t, account.Validate())

	// Invalid Account Number is provided
	account.AccountData.Attributes.AccountNumber = "5u426819"
	assert.Equal(t, "Key: 'Account.AccountData.Attributes.AccountNumber' Error:Field validation for 'AccountNumber' failed on the 'number' tag", account.Validate().Error())

	// Invalid Bank ID Code is provided
	account.AccountData.Attributes.AccountNumber = "51426819"
	account.AccountData.Attributes.BankIDCode = "dDHUF"
	assert.Equal(t, "Key: 'Account.AccountData.Attributes.BankIDCode' Error:Field validation for 'BankIDCode' failed on the 'bank_id_code' tag", account.Validate().Error())

	// Invalid BIC is provided
	account.AccountData.Attributes.BankIDCode = "RDHUF"
	account.AccountData.Attributes.BIC = "iWBKGB22"
	assert.Equal(t, "Key: 'Account.AccountData.Attributes.BIC' Error:Field validation for 'BIC' failed on the 'bic' tag", account.Validate().Error())

	// Invalid IBAN is provided
	account.AccountData.Attributes.BIC = "IWBKGB22"
	account.AccountData.Attributes.IBAN = "Gi11NWBK40030041426819"
	assert.Equal(t, "Key: 'Account.AccountData.Attributes.IBAN' Error:Field validation for 'IBAN' failed on the 'iban' tag", account.Validate().Error())

	// Invalid UUID is provided for OrganisationID
	account.AccountData.Attributes.IBAN = "GL11NWBK40030041426819"
	account.AccountData.OrganisationID = "b0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	assert.Equal(t, "Key: 'Account.AccountData.OrganisationID' Error:Field validation for 'OrganisationID' failed on the 'oid' tag", account.Validate().Error())

	// All Valid informations
	account.AccountData.ID = "9127e265-9605-4b4b-a0e5-3003ea9cc4dc"
	account.AccountData.OrganisationID = "2b0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	assert.Nil(t, account.Validate())

	// Invalid Long Title is provided
	account.AccountData.Title = "This string is for testing which will take error because it will be longer than 40 characters"
	assert.Equal(t, "Key: 'Account.AccountData.Title' Error:Field validation for 'Title' failed on the 'title' tag", account.Validate().Error())
}
