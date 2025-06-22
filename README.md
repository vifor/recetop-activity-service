# Recetop Activity Service

A high-performance microservice built with Go (Golang) for tracking and ingesting user activity events from the Recetop application. This service is designed to be lightweight, efficient, and scalable, making it ideal for a serverless deployment on AWS Lambda.

This project follows a **"Design-First" API approach**, with the API contract defined in an OpenAPI 3.0 specification. The Go code for the server interface and data models is generated from this specification, ensuring the implementation always matches the contract.

---

## Technologies Used

- **Language:** Go (Golang)
- **API Specification:** OpenAPI 3.0
- **Code Generation:** oapi-codegen
- **Web Framework:** Echo v4
- **Target Platform:** AWS Lambda (using a `provided.al2023` custom runtime)
- **Lambda Adapter:** algnhsa
- **Infrastructure:** AWS API Gateway (HTTP API), AWS IAM

---

## API Contract and Code Generation

This project uses **oapi-codegen** to generate Go code from the API specification. The `openapi.yaml` file is the single source of truth for the API's contract.

If you make changes to `openapi.yaml`, you must regenerate the server code before building. The generation command is embedded in the `generate.go` file.

To regenerate the API code, run the following command from the project root:

```bash
go generate ./...
```

---

## Deployment to AWS Lambda

### Prerequisites

- Go (version 1.21 or later recommended)
- An AWS account with access to IAM, Lambda, and API Gateway.

---

### Step 1: Build the Executable

Compile the application into a Linux executable named `bootstrap`. This name is required by the AWS Lambda custom runtime.

**On Windows (PowerShell):**

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bootstrap .
```

**On macOS / Linux (Bash/Zsh):**

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap .
```

---

### Step 2: Package for Deployment

Compress the `bootstrap` binary into a `.zip` file for uploading to Lambda.

**On Windows (PowerShell):**

```powershell
Compress-Archive -Path .\bootstrap -DestinationPath .\deployment.zip
```

**On macOS / Linux (Bash/Zsh):**

```bash
zip deployment.zip bootstrap
```

---

### Step 3: Deploy to AWS

Upload the `deployment.zip` package to a new or existing AWS Lambda function with the following configuration:

- **IAM Role:** Requires the `AWSLambdaBasicExecutionRole` policy to write logs to CloudWatch.
- **Lambda Function:**
  - **Runtime:** Custom runtime on Amazon Linux 2023 (`provided.al2023`).
  - **Handler:** `bootstrap`.
  - **Code Upload:** Upload your `deployment.zip` file.
- **API Gateway Trigger:**
  - Add an HTTP API trigger to make the function accessible via a public URL.

---

## Local Testing (Advanced)

While the primary target is AWS Lambda, the Hexagonal Architecture makes local testing simple. For advanced local testing that fully mimics the Lambda environment, the [AWS Serverless Application Model (SAM) CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli.html) is recommended.

---

## API Endpoint

*Note: Due to the default AWS HTTP API Gateway configuration, all routes are prefixed with the stage and function name. The application's router is configured to handle this base path (`/default/recetop-activity-service`).*

| HTTP Method | Endpoint | Description                                | Response Body                                              |
| :---------- | :------- | :----------------------------------------- | :--------------------------------------------------------- |
| `GET`       | `/`      | Checks the health of the service endpoint. | `{"status":"UP","message":"Go serverless with OpenAPI spec!"}` |
---

## Testing Strategy

This project uses a multi-layered testing approach to ensure code quality, correctness, and reliability. The strategy is inspired by the "testing pyramid," with a focus on fast, automated tests that run within the Go toolchain.

### Unit Tests

* **Tool:** Go's standard `testing` package.
* **Goal:** To test the core application logic (the "Hexagon" in our architecture) in complete isolation. This involves testing the methods on the `Server` struct directly, without any dependency on the Echo framework or AWS infrastructure. These tests are extremely fast and verify the business logic.

### Integration Tests (API Layer)

* **Tool:** Go's `testing` and `net/http/httptest` packages.
* **Goal:** To test the full application stack (HTTP routing, middleware, handlers, and core logic) without deploying to AWS. These tests spin up an in-memory version of the Echo server, make real HTTP requests to it, and validate the HTTP responses (status codes, headers, and JSON bodies). This provides a fast and reliable way to verify that the entire API is functioning correctly before deployment.

### End-to-End (E2E) Tests

* **Tool:** Postman / cURL.
* **Goal:** To verify that the service is working correctly after being deployed to the live AWS environment. This is a final sanity check to ensure the infrastructure (API Gateway, Lambda permissions) is configured correctly.

### Architectural Decision: `httptest` vs. Karate

For the integration/API testing layer, I've considered external, black-box testing tools like [Karate](https://github.com/karatelabs/karate), which is excellent for testing any HTTP endpoint regardless of the backend language.

However, I've made the deliberate decision to use Go's native `net/http/httptest` library for this project's primary integration tests. The reasoning is as follows:

* **Tight Integration:** The tests are written in Go and run with the standard `go test` command. There is no need for external runtimes (like a JVM) or dependencies.
* **Extreme Speed:** In-memory tests execute in milliseconds, providing an incredibly fast feedback loop for developers. This encourages running tests frequently.
* **No Deployment Necessary:** The full application stack can be validated without the time-consuming process of deploying to a remote environment.
* **Compile-Time Safety:** The tests are part of the Go source code and benefit from the same type-safety and compile-time checks as the application itself.

While Karate is a fantastic tool for true E2E testing against a deployed environment (often managed by a separate QA team), the `httptest` approach was chosen to maximize developer velocity and create a robust, self-contained safety net within the project itself.

## Architectural Decisions

### Choice of Language: Go

Go was deliberately chosen to optimize for performance and cost in a cloud environment. Its fast compile times, small binary size, low memory footprint, and efficient concurrency result in significantly faster "cold starts" on AWS Lambda compared to interpreted languages.

---

### Design Pattern: Hexagonal Architecture (Ports and Adapters)

This service is structured following the Hexagonal Architecture pattern to isolate the core application logic from external concerns like AWS Lambda.

- **The Hexagon (Core Logic):** The `Server` struct contains the pure business logic for handling API requests. It has no knowledge of AWS.
- **The Port (The API Contract):** The `api.ServerInterface` generated by oapi-codegen from `openapi.yaml`. This Go interface is the formal contract that the core logic must implement.
- **The Adapter:** The `main` function acts as the adapter. It uses the `algnhsa` library to translate incoming AWS Lambda / API Gateway events into standard HTTP requests that are then routed to the core application via the Echo framework.

This separation provides immense benefits, including enhanced testability (the core logic can be unit-tested without AWS) and portability (the core logic could be served by a standard `net/http` server with a different adapter).

---

### HTTP Handling: Echo Framework

Instead of handling raw HTTP requests, the service uses the Echo web framework.

- **Performance:** Echo is a high-performance, minimalist framework that adds very little overhead.
- **Productivity:** It provides helpful utilities for routing, data binding, and response generation, which simplifies the application code. Its middleware system is useful for cross-cutting concerns like logging and error recovery.
- **Integration:** oapi-codegen integrates seamlessly with Echo, generating compatible server interfaces out of the box.