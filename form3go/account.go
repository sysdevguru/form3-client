package form3go

import (
	"regexp"

	"github.com/pariz/gountries"
	"golang.org/x/text/currency"
	"gopkg.in/go-playground/validator.v9"
)

// Account represents Form3 Organisation Account
type Account struct {
	AccountData Data `json:"data"`
}

// Data is account information
type Data struct {
	Type                        string            `json:"type" validate:"type"`
	ID                          string            `json:"id" validate:"id"`
	OrganisationID              string            `json:"organisation_id" validate:"oid"`
	Attributes                  AccountAttributes `json:"attributes"`
	Title                       string            `json:"title" validate:"title"`
	FirstName                   string            `json:"first_name" validate:"first_name"`
	BankAccountName             string            `json:"bank_account_name" validate:"ban"`
	AlternativeBankAccountNames []string          `json:"alternative_bank_account_names"`
	AccountClassification       string            `json:"account_classification"`
	JointAccount                bool              `json:"joint_account"`
	AccountMatchingOptOut       bool              `json:"account_matching_opt_out"`
	SecondaryIdentification     string            `json:"secondary_identification" validate:"si"`
}

// AccountAttributes is Account Attributes
type AccountAttributes struct {
	Country       string `json:"country" validate:"country"`
	BaseCurrency  string `json:"base_currency" validate:"currency"`
	AccountNumber string `json:"account_number" validate:"number"`
	BankID        string `json:"bank_id" validate:"bank_id"`
	BankIDCode    string `json:"bank_id_code" validate:"bank_id_code"`
	BIC           string `json:"bic" validate:"bic"`
	IBAN          string `json:"iban" validate:"iban"`
}

// Validate validates Account fields
func (a Account) Validate() error {
	v := validator.New()
	_ = v.RegisterValidation("country", func(fl validator.FieldLevel) bool {
		query := gountries.New()
		code := fl.Field().String()
		_, err := query.FindCountryByAlpha(code)
		if err != nil {
			return false
		}
		return true
	})
	_ = v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		cur := fl.Field().String()
		if cur != "" {
			_, err := currency.ParseISO(cur)
			if err != nil {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("number", func(fl validator.FieldLevel) bool {
		number := fl.Field().String()
		if number != "" {
			rxPat := regexp.MustCompile("^[A-Z0-9]{0,64}$")
			if !rxPat.MatchString(number) {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("bank_id", func(fl validator.FieldLevel) bool {
		bankID := fl.Field().String()
		if bankID != "" {
			rxPat := regexp.MustCompile("^[A-Z0-9]{0,16}$")
			if !rxPat.MatchString(bankID) {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("bank_id_code", func(fl validator.FieldLevel) bool {
		bankIDCode := fl.Field().String()
		if bankIDCode != "" {
			rxPat := regexp.MustCompile("^[A-Z]{0,16}$")
			if !rxPat.MatchString(bankIDCode) {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("bic", func(fl validator.FieldLevel) bool {
		bic := fl.Field().String()
		if bic != "" {
			rxPat := regexp.MustCompile("^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$")
			if !rxPat.MatchString(bic) {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("iban", func(fl validator.FieldLevel) bool {
		iban := fl.Field().String()
		if iban != "" {
			rxPat := regexp.MustCompile("^[A-Z]{2}[0-9]{2}[A-Z0-9]{0,64}$")
			if !rxPat.MatchString(iban) {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("title", func(fl validator.FieldLevel) bool {
		title := fl.Field().String()
		if title != "" {
			if len(title) > 40 {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("first_name", func(fl validator.FieldLevel) bool {
		first_name := fl.Field().String()
		if first_name != "" {
			if len(first_name) > 40 {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("ban", func(fl validator.FieldLevel) bool {
		ban := fl.Field().String()
		if ban != "" {
			if len(ban) > 140 {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("si", func(fl validator.FieldLevel) bool {
		secondID := fl.Field().String()
		if secondID != "" {
			if len(secondID) > 140 {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("type", func(fl validator.FieldLevel) bool {
		typeName := fl.Field().String()
		if typeName != "accounts" {
			return false
		}
		return true
	})
	_ = v.RegisterValidation("id", func(fl validator.FieldLevel) bool {
		id := fl.Field().String()
		r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		if r.MatchString(id) {
			return true
		}
		return false
	})
	_ = v.RegisterValidation("oid", func(fl validator.FieldLevel) bool {
		oid := fl.Field().String()
		if oid != "" {
			r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
			if !r.MatchString(oid) {
				return false
			}
		}
		return true
	})

	return v.Struct(a)
}
