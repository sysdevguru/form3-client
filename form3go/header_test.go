package form3go

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDigestHeader(t *testing.T) {
	requestInfo := request{
		data: `{"data":{"type":"payments","id":"1234567890","version":0,"organisation_id":"1234567890","attributes":{"amount":"200.00","beneficiary_party":{"account_name":"Mrs Receiving Test","account_number":"71268996","account_number_code":"BBAN","account_with":{"bank_id":"400302","bank_id_code":"GBDSC"}},"currency":"GBP","debtor_party":{"account_name":"Mr Sending Test","account_number":"87654321","account_number_code":"BBAN","account_with":{"bank_id":"1234567890","bank_id_code":"GBDSC"}},"processing_date":"2019-20-5","reference":"Something","payment_scheme":"FPS","scheme_payment_sub_type":"TelephoneBanking","scheme_payment_type":"ImmediatePayment"}}}`,
	}

	assert.Equal(t, "SHA-256=WllU95a/P37KDBmTedpEIIvVtBgRqDdYrHz06NXDuvk=", requestInfo.genDigestHeader())
}

func TestGenSignature(t *testing.T) {
	requestInfo := request{
		data:     testAccountInfo,
		method:   "POST",
		endpoint: "/v1/organisation/accounts",
		keyID:    os.Getenv("FORM3_KEY_ID"),
		keyPath:  os.Getenv("FORM3_PRIV_KEY_PATH"),
	}
	date := "Wed, 08 Jan 2020 03:52:44 EST"
	digest := requestInfo.genDigestHeader()
	sig, _ := requestInfo.genSignature(date, digest)
	assert.Equal(t, "tCLWrCm2Tm2EQ0uoIKmai0cifub+R5XI/huAul1IVWHqNBifgvTFtJc1HT1nNtv85PZzpgdzTqAuZOd1jn1d8EXEQbjC/79koiE1z9REuMdKxIwk5KKmZVRQkaHUAg51mO8uAZRYJrpOnGUROoG4hU4WC0fOwNeBW36XIrpVlYM=", sig)

	// Invalid Key_ID is provided
	requestInfo.keyID = ""
	sig, err := requestInfo.genSignature(date, digest)
	assert.Equal(t, "", sig)
	assert.Equal(t, "empty FORM3_KEY_ID env variable", err.Error())

	// Invalid Key Path is provided
	requestInfo.keyID = os.Getenv("FORM3_KEY_ID")
	requestInfo.keyPath = ""
	sig, err = requestInfo.genSignature(date, digest)
	assert.Equal(t, "", sig)
	assert.Equal(t, "empty FORM3_PRIV_KEY_PATH env variable", err.Error())
}

func TestGenAuthHeader(t *testing.T) {
	requestInfo := request{
		data:     testAccountInfo,
		method:   "POST",
		endpoint: "/v1/organisation/accounts",
		keyID:    os.Getenv("FORM3_KEY_ID"),
		keyPath:  os.Getenv("FORM3_PRIV_KEY_PATH"),
	}
	date := "Wed, 08 Jan 2020 03:52:44 EST"
	digest := requestInfo.genDigestHeader()
	sig, _ := requestInfo.genSignature(date, digest)
	authHeader, _ := requestInfo.genAuthHeader(sig)
	expectedHeader := `Signature keyId="75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8",algorithm="rsa-sha256",header="(request-target) host date accept content-type content-length digest",signature="tCLWrCm2Tm2EQ0uoIKmai0cifub+R5XI/huAul1IVWHqNBifgvTFtJc1HT1nNtv85PZzpgdzTqAuZOd1jn1d8EXEQbjC/79koiE1z9REuMdKxIwk5KKmZVRQkaHUAg51mO8uAZRYJrpOnGUROoG4hU4WC0fOwNeBW36XIrpVlYM="`
	assert.Equal(t, expectedHeader, authHeader)
}
