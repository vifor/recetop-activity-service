# Recetop Activity Service

A high-performance microservice built with Go (Golang) for tracking and ingesting user activity events from the Recetop application. This service is designed to be lightweight, efficient, and scalable, making it ideal for a serverless deployment on AWS Lambda.

---

## Technologies Used

- **Language:** Go (Golang)
- **Target Platform:** AWS Lambda (using a `provided.al2023` custom runtime)
- **Infrastructure:** AWS API Gateway (HTTP API), AWS IAM

---

## Deployment to AWS Lambda

This service is not intended to be run as a traditional local server. It is designed to be compiled and deployed directly to AWS Lambda.

### Prerequisites

- Go (version 1.20 or later recommended)
- An AWS account with access to IAM, Lambda, and API Gateway.

---

### Step 1: Build the Executable

To deploy this function, you must first compile it into a Linux executable named `bootstrap`. This specific name is a requirement for AWS Lambda's custom runtimes (`provided.al2023`).

**On Windows (PowerShell):**

```powershell
# Set environment variables to compile for Linux, then build the 'bootstrap' file.
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bootstrap .
```

**On macOS / Linux (Bash/Zsh):**

```bash
# Set environment variables to compile for Linux, then build the 'bootstrap' file.
GOOS=linux GOARCH=amd64 go build -o bootstrap .
```

---

### Step 2: Package for Deployment

AWS Lambda expects the code to be in a `.zip` file. Compress the `bootstrap` binary you just created.

**On Windows (PowerShell):**

```powershell
# Create a zip archive containing the bootstrap executable.
Compress-Archive -Path .\bootstrap -DestinationPath .\deployment.zip
```

**On macOS / Linux (Bash/Zsh):**

```bash
# Create a zip archive containing the bootstrap executable.
zip deployment.zip bootstrap
```

---

### Step 3: Deploy to AWS

Upload the `deployment.zip` package to a new AWS Lambda function with the following configuration:

- **IAM Role:** Create a new IAM Role for the function with the `AWSLambdaBasicExecutionRole` managed policy. This allows the function to write logs to CloudWatch.
- **Lambda Function:**
  - **Runtime:** Select Custom runtime on Amazon Linux 2023 (`provided.al2023`).
  - **Handler:** Set the handler name to `bootstrap` (or `main`, as the bootstrap file in the zip takes precedence).
  - **Code Upload:** Upload your `deployment.zip` file.
- **API Gateway Trigger:**
  - Add a trigger to your Lambda function.
  - Select API Gateway.
  - Choose to create a new HTTP API.
  - This will provide you with a public URL to invoke your function.

---

## Local Testing (Advanced)

Directly running this application with `go run` is not possible as it lacks a web server. For local testing that mimics the cloud environment, the recommended approach is to use the [AWS Serverless Application Model (SAM) CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli.html). This tool can invoke your function locally in a Docker container that replicates the Lambda environment.

---

## API Endpoint

| HTTP Method | Endpoint | Description                          | Response Body                                      |
|-------------|----------|--------------------------------------|----------------------------------------------------|
| GET         | /        | Checks the health of the service endpoint. | `{"status":"UP","message":"Go serverless with AWS Lambda!"}` |

---

## Architectural Decisions

### Choice of Language: Go

Go was deliberately chosen for this serverless microservice to optimize for performance and cost in a cloud environment.

- **Reduced Cold Start Times:** As a compiled language, Go produces a small, self-contained binary. This allows the Lambda function to have significantly faster "cold start" times compared to interpreted languages, reducing latency and improving user experience.
- **Efficiency:** Go's low memory footprint and efficient CPU usage are ideal for a serverless platform where allocated resources and execution time are directly tied to cost.
- **Modern Custom Runtime:** By compiling for `provided.al2023`, we ensure the application runs in the latest, most secure, and most performant Lambda environment, giving us full control over the