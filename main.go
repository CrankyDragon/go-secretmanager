package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {

	var (
		extract = flag.Bool("e", true, "Extract json")
		secret  = flag.String("s", "secret", "Secret To Fetch")
		region  = flag.String("r", "us-east-1", "AWS Region")
		version = flag.String("v", "version", "Version of secret To Fetch")
	)

	flag.Parse()
	if *secret == "secret" {
		fmt.Println("You must specify a secret name to fetch")
		return
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: region},
	}))

	getSecret(sess, secret, version, extract)
}

func getSecret(sess *session.Session, secretName, secretVersion *string, extract *bool) {
	svc := secretsmanager.New(sess)
	var versionID string
	if *secretVersion == "version" {
		versionID = "AWSCURRENT"
	} else {
		versionID = *secretVersion
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(*secretName),
		VersionStage: aws.String(versionID),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	if *extract {
		fmt.Println(*result.SecretString)
	} else {
		data, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
	}
}
