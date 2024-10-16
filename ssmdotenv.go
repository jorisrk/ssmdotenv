package ssmdotenv

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

var verbose bool
var prefix string
var client SSMClient

// SetVerbose sets the verbosity level. If set to true, verbose logs will be enabled.
func SetVerbose(v bool) {
	verbose = v
}

// SetPrefix sets a prefix to be added to SSM parameter names when fetching them.
func SetPrefix(p string) {
	prefix = p
}

// Load loads environment variables from a .env file and optionally from AWS SSM parameters.
// It accepts one or more paths to AWS SSM, and retrieves parameters from each path.
func Load(paths ...string) {
	err := godotenv.Load()
	if err != nil {
		verboseLog(".env file not found")
	}

	client := getSsmClient()

	if client == nil {
		return
	}

	for _, path := range paths {
		paginator := client.GetParametersByPathPaginator(&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		})

		count := 0

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.Background())
			if err != nil {
				verboseLog("Error getting parameters:", err)
				return
			}

			for _, param := range page.Parameters {
				name := (*param.Name)[len(path):]
				if os.Getenv(name) != "" {
					continue
				}
				verboseLog("Loaded parameter:", name)
				os.Setenv(name, *param.Value)
				count++
			}
		}

		if count == 0 {
			verboseLog("No parameters loaded from SSM for path", path)
		}
	}
}

// Env retrieves the value of an environment variable by key.
// If the variable is not found, it returns the provided defaultValue or an empty string if no defaultValue is specified.
func Env(key string, defaultValue ...string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return def(defaultValue...)
}

// GetParameter fetches a single parameter from AWS SSM.
// It uses the prefix set by SetPrefix() and appends it to the parameterName.
// Returns the parameter value, or defaultValue if the parameter is not found.
func GetParameter(parameterName string, defaultValue ...string) string {
	client := getSsmClient()

	if client == nil {
		return def(defaultValue...)
	}

	path := prefix + parameterName

	// Retrieve the parameter
	input := &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	}

	result, err := client.GetParameter(context.Background(), input)
	if err != nil {
		verboseLog("Unable to get parameter:", path, err)
		return def(defaultValue...)
	}

	return *result.Parameter.Value
}

// def is a helper function that returns the first value from defaultValue, or an empty string if none is provided.
func def(defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// getSsmClient initializes and returns an AWS SSM client.
// If the client has already been created, it returns the existing client.
func getSsmClient() SSMClient {
	if client != nil {
		return client
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(Env("AWS_REGION", "eu-west-3")),
	)

	if err != nil {
		verboseLog("Can't create AWS SSM client: ", err)
		return nil
	}

	client = &AWSSSMClient{client: ssm.NewFromConfig(cfg)}
	return client
}

// verboseLog logs messages only when verbose mode is enabled.
func verboseLog(v ...any) {
	if !verbose {
		return
	}

	log.Println(v...)
}
