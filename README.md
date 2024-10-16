# ssmdotenv

`ssmdotenv` is a Go package that simplifies loading environment variables from both local `.env` files and AWS Systems Manager (SSM) Parameter Store. It provides a convenient way to manage environment variables, especially when dealing with multiple environments (e.g., local development, staging, production).

## Features

-   Load environment variables from a `.env` file.
-   Fetch parameters from AWS SSM Parameter Store, supporting recursive and decrypted parameters.
-   Set custom prefixes for SSM parameters to avoid conflicts.
-   Verbose mode for detailed logging.

## Installation

To install the package, use `go get`:

```bash
go get github.com/jorisrk/ssmdotenv
```

## Usage

### 1. Loading Environment Variables from `.env` File

The package can read variables from a local `.env` file.

```go
package main

import (
    "github.com/jorisrk/ssmdotenv"
)

func main() {
    // Enable verbose logging
    ssmdotenv.SetVerbose(true)

    // Load environment variables from .env
    ssmdotenv.Load()

    // Access the environment variable
    value := ssmdotenv.Env("MY_ENV_VAR", "default_value")
    println(value)
}
```

### 2. Loading Parameters from AWS SSM

You can also load environment variables directly from AWS SSM Parameter Store. Make sure your AWS credentials are configured.

```go
package main

import (
    "github.com/jorisrk/ssmdotenv"
)

func main() {
    // Enable verbose logging
    ssmdotenv.SetVerbose(true)

    // Set a custom prefix for SSM parameters (optional)
    ssmdotenv.SetPrefix("/myapp/")

    // Load environment variables from AWS SSM Parameter Store
    ssmdotenv.Load("/myapp/")

    // Retrieve a parameter from SSM
    value := ssmdotenv.GetParameter("MY_SSM_PARAM", "default_value")
    println(value)
}
```

### 3. Combining `.env` File and SSM

You can combine both approaches. The package will prioritize environment variables loaded from `.env` over SSM if there are duplicates.

```go
package main

import (
    "github.com/jorisrk/ssmdotenv"
)

func main() {
    ssmdotenv.SetVerbose(true)

    // Load both .env and SSM parameters
    ssmdotenv.Load("/myapp/")
    value := ssmdotenv.Env("MY_ENV_VAR", "default_value")
    println(value)
}
```

## Environment Variables

-   `AWS_REGION`: Set this environment variable to specify the AWS region for SSM. Defaults to `eu-west-3` if not specified.

## Verbose Mode

Verbose mode is useful for debugging. It will print logs about which environment variables and SSM parameters are being loaded.

```go
ssmdotenv.SetVerbose(true)
```
