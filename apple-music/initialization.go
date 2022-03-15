package apple_music

// Retrieves Apple Music Token
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"log"
)

type Secret struct {
	Key   string
	Value string
}

func getSecret() string {
	secretName := "FujiAppleMusicToken"
	region := "us-east-1"

	// Create a Secrets Manager client
	// The credentials file needs to be setup in your .aws directory in your home directory
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})
	if err != nil {
		// Handle session creation error
		log.Println("Unable to create a New Session while retrieving a secret.")
		log.Println(err.Error())
		return ""
	}

	// Setup configuration with region and credentials
	// This ensures credentials are setup
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		// Handle session creation error
		log.Println("Error retrieving credentials.")
		log.Println(err.Error())
		return ""
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.Println("Unable to retrieve secrets value in AWS")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				log.Println("Secrets Manager can't decrypt the protected secret text using the provided KMS key.")
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				log.Println("An error occurred on the server side.")
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				log.Println("You provided an invalid value for a parameter.")
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				log.Println("You provided a parameter value that is not valid for the current state of the resource.")
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				log.Println("We can't find the resource that you asked for.")
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println("Unknown error occurred while retrieving secret " + secretName)
			log.Println(err.Error())
		}
		return ""
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// Removing binary secret as not used at the moment
	var secretStringJson, _ string
	if result.SecretString != nil {
		secretStringJson = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Println("Base64 Decode Error:", err)
			return ""
		}
		_ = string(decodedBinarySecretBytes[:len])
	}

	// Parse key/value string
	secretMap := map[string]string{}
	json.Unmarshal([]byte(secretStringJson), &secretMap)

	// Return secret
	return secretMap[secretName]
}
