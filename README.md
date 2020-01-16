# Form3 Account API Go client

Simple Golang client library for Account API  

## Prerequisites
| Environment variable | Description                                |
|:---------------------|:-------------------------------------------|
| FORM3_HOST           | AccountAPI URL                             |
| FORM3_KEY_ID         | Public Key ID                              |
| FORM3_PRIV_KEY_PATH  | Private Key Path                           |

### Create form3go client
```go
client := form3go.Client{
    PubKeyID:    os.Getenv("FORM3_KEY_ID"),
    PrivKeyPath: os.Getenv("FORM3_PRIV_KEY_PATH"),
    HttpClient:  http.Client{},
}
```

### Create Account
CreateAccount returns account if creating account succeed.
```go
account := form3go.Account{
    // account informations
}
acct, err := client.CreateAccount(account)
```

### Fetch Account
```go
id := "Account ID here"
acct, err := client.FetchAccount(id)
```

### List Accounts
```go
accounts, _ := client.ListAccounts(pageNumber, pageSize) // pageNumber, pageSize are int values
```

### Delete Account
```go
id := "Account ID here"
err := client.DeleteAccount(id)
```
