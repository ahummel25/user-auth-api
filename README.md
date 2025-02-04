# User Auth API

A GraphQL-based user authentication API built with Go and the Serverless Framework, designed to run on AWS Lambda.

## Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Serverless Framework](https://www.serverless.com/framework/docs/getting-started)
- AWS Account and [configured credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

## Local Development Setup

1. Clone the repository:
   ```sh
   git clone git@github.com:ahummel25/user-auth-api.git
   cd user-auth-api
   ```

2. Install dependencies:
   ```sh
   go mod download
   go mod tidy
   npm install -g serverless
   npm install
   ```

3. Set up local environment variables:
   ```sh
   # Copy the example env file and modify as needed
   cp .env.example .env
   ```

4. Start local MongoDB:
   ```sh
   make dev-deps
   ```

5. Run the development server:
   ```sh
   make local
   ```

The API will be available at:
- GraphQL API: http://localhost:8080/graphql
- GraphiQL Playground: http://localhost:8080/graphiql
- Apollo Playground: http://localhost:8080/apollo

## Development Commands

```sh
# Start MongoDB
make dev-deps

# Run local server
make local

# View MongoDB logs
make logs-mongo

# Run tests
go test ./...

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── cmd/                    # Command line tools
│   └── local/             # Local development server
├── config/                # Configuration management
├── db/                    # Database layer
├── graphql/              # GraphQL schema and resolvers
│   ├── directives/       # GraphQL directives
│   ├── generated/        # Generated GraphQL code
│   ├── model/           # GraphQL models
│   ├── resolvers/       # GraphQL resolvers
│   └── schema/          # GraphQL schema definitions
├── lambda/               # AWS Lambda functions
│   └── graphql/         # GraphQL API Lambda
├── service/             # Business logic services
└── utils/               # Utility functions
```

## Deployment

The API is deployed using the Serverless Framework to AWS Lambda.

1. Configure AWS credentials:
   ```sh
   aws configure
   ```

2. Deploy to AWS:
   ```sh
   # Deploy to dev
   export TF_WORKSPACE=dev
   make deploy

   # Deploy to prod
   export TF_WORKSPACE=prod
   make deploy
   ```

## Available Environments

- **Local**: Local development environment
- **Dev**: Development environment in AWS
- **Prod**: Production environment in AWS

## Monitoring and Logs

### Local Development
- API logs are visible in the terminal running `make local`
- MongoDB logs can be viewed with `make logs-mongo`

### AWS Environments
- CloudWatch Logs
- X-Ray Tracing (enabled for API Gateway and Lambda)
