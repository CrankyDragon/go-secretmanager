Go Based AWS SecretManager

Simple command line for fetching secrets from AWS Secrets Manager

#### Command line Arguments
Currently supported
```
  -e bool
        Output's the SecretString value.
  -s string
        Secret To Fetch (default "secret")
  -r string
        AWS Region (default "us-east-1")
  -v string
        Version of secret To Fetch (default "version")
```