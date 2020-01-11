## Client prerequisites
Clients have to export env variables in order to create client object of this library.

| Environment variable | Description                                |
|:---------------------|:-------------------------------------------|
| FORM3_HOST           | AccountAPI url, localhost:8080             |
| FORM3_KEY_ID         | Public Key ID                              |
| FORM3_PRIV_KEY_PATH  | Private Key Path                           |

I used FORM3 test_private_key.pem and public key_id for this task.
Those env variables are defined in `common.env` file.

## Authorization
I assume that all the clients will use message signing for the authentication and validation.
Deprecated Basic authentication will not be handled.

## Account struct validation
I validate provided Account struct when creating account.
Even though provided account_api validates the struct.

## Running task
Sometimes the containers are not running as described in the `depends_on`
The reason could be existing, non removed containers
So I recommend to run following command
```sh
docker-compose up --build --remove-orphans
```

service `test` is my implementation and it will run unit & integration tests